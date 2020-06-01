/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"

	"github.com/go-trellis/config"
	"github.com/spf13/cobra"
)

// yearCmd represents the year command
var yearCmd = &cobra.Command{
	Use:   "year",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		yt := &YearTax{}
		r := config.NewYAMLReader()
		if err := r.Read(cfgFile, yt); err != nil {
			panic(err)
		}
		yt.run()
	},
}

func init() {
	rootCmd.AddCommand(yearCmd)
}

type YearTax struct {
	Global YearTaxGlobalConfig `yaml:"global"`

	Salaries     []*MonthSalary `yaml:"salaries"`
	YearTaxRates []*YearTaxRate `yaml:"year_tax_rates"`

	totalSalaries    float64 `yaml:"-"`
	totalTaxSalaries float64 `yaml:"-"`
	totalTaxation    float64 `yaml:"-"`
}

type YearTaxGlobalConfig struct {
	For                bool `yaml:"for"`
	InsurancesAndFound `yaml:",inline"`
}

type MonthSalary struct {
	Threshold        float64 `yaml:"threshold"`         // 起征点
	Salary           float64 `yaml:"salary"`            // 薪水
	SubsidyAmount    float64 `yaml:"subsidy_amount"`    // 补贴
	Insurance        float64 `yaml:"insurance"`         // 保险金额
	AccumulationFund float64 `yaml:"accumulation_fund"` // 公积金
	DeductibleAmount float64 `yaml:"deductible_amount"` // 抵扣金额
}

type YearTaxRate struct {
	SalaryMin      float64 `yaml:"salary_min"`
	SalaryMax      float64 `yaml:"salary_max"`
	Rate           float64 `yaml:"rate"`
	DeductedAmount float64 `yaml:"deducted_amount"`
}

type MonthlyTax struct {
	MonthSalary     *MonthSalary `yaml:"month_salary"`
	Month           int          `yaml:"month"`
	Taxation        float64      `yaml:"taxation"`
	Salary          float64      `yaml:"salary"`
	HistoryTaxation float64      `yaml:"history_taxation"`
	HistorySalary   float64      `yaml:"history_salary"`
}

type InsurancesAndFound struct {
	AccumulationFundBase float64 `yaml:"accumulation_fund_base"`
	AccumulationFundRate float64 `yaml:"accumulation_fund_rate"`

	InsurancesBase                float64 `yaml:"insurances_base"`
	EndowmentInsuranceRate        float64 `yaml:"endowment_insurance_rate"`
	UnemploymentInsuranceRate     float64 `yaml:"unemployment_insurance_rate"`
	EmploymentInjuryInsuranceRate float64 `yaml:"employment_injury_insurance_rate"`
	BirthInsuranceRate            float64 `yaml:"birth_insurance_rate"`
	MedicalRate                   float64 `yaml:"medical_rate"`
}

func (p *YearTax) run() {
	var taxes []*MonthlyTax
	var lastMonth int
	for i, monthSalary := range p.Salaries {
		lastMonth = i + 1
		taxes = append(taxes, p.getMonthTax(lastMonth, monthSalary))
	}

	if p.Global.For {
		for i := lastMonth + 1; i < 13; i++ {
			monthSalary := p.Salaries[lastMonth-1]
			p.Salaries = append(p.Salaries, monthSalary)
			taxes = append(taxes, p.getMonthTax(i, monthSalary))
		}
	}

	for _, tax := range taxes {
		fmt.Printf("%0.2d月，个税: %0.2f，总个税：%0.2f，薪资: %0.2f，总薪资: %0.2f\n",
			tax.Month, tax.Taxation, tax.HistoryTaxation, tax.Salary, tax.HistorySalary)
	}
}

func (p *YearTax) getMonthTax(lastMonth int, monthSalary *MonthSalary) *MonthlyTax {

	if monthSalary.AccumulationFund == 0 {
		monthSalary.AccumulationFund = p.Global.AccumulationFundBase * p.Global.AccumulationFundRate / 100.0
	}

	if monthSalary.Insurance == 0 {
		monthSalary.Insurance = p.Global.InsurancesBase *
			(p.Global.AccumulationFundRate + p.Global.UnemploymentInsuranceRate + p.Global.MedicalRate +
				p.Global.EmploymentInjuryInsuranceRate + p.Global.EndowmentInsuranceRate) /
			100.0
	}

	salary := monthSalary.Salary + monthSalary.SubsidyAmount - monthSalary.Insurance - monthSalary.AccumulationFund
	taxSalary := salary - monthSalary.Threshold - monthSalary.DeductibleAmount
	if taxSalary > 0 {
		p.totalTaxSalaries += taxSalary
	}
	p.totalSalaries += salary

	// 小于起征点，那么税收为0
	if salary-monthSalary.DeductibleAmount <= monthSalary.Threshold {
		return &MonthlyTax{
			MonthSalary: monthSalary,
			Month:       lastMonth,
			Taxation:    0,
			Salary:      salary,
		}
	}

	for _, taxRate := range p.YearTaxRates {
		if p.totalTaxSalaries <= taxRate.SalaryMin ||
			(p.totalTaxSalaries > taxRate.SalaryMax && taxRate.SalaryMax != 0) {
			continue
		}
		tax := (p.totalTaxSalaries*taxRate.Rate)/100.0 - taxRate.DeductedAmount - p.totalTaxation
		p.totalTaxation += tax

		return &MonthlyTax{
			MonthSalary:     monthSalary,
			HistoryTaxation: p.totalTaxation,
			HistorySalary:   p.totalSalaries,
			Month:           lastMonth,
			Taxation:        tax,
			Salary:          salary - tax,
		}
	}
	return nil
}
