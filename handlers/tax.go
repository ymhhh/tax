/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package handlers

import (
	"fmt"

	"github.com/go-trellis/config"
)

// TaxesHandler handler 对象
type TaxesHandler struct {
	AccumulationFundHandler `yaml:",inline" json:",inline"`
	InsurancesHandler       `yaml:",inline" json:",inline"`

	YearTaxBase `yaml:",inline" json:",inline"`

	totalSalaries    float64
	totalTaxSalaries float64
	totalTaxation    float64
}

// NewTaxesHandler 生成handler对象
func NewTaxesHandler(file string) (*TaxesHandler, error) {
	t := &TaxesHandler{}
	if err := config.NewSuffixReader().Read(file, t); err != nil {
		return nil, err
	}
	return t, nil
}

// Salaries 薪资配置参数
type Salaries struct {
	For bool `yaml:"for" json:"for"`

	PersonalInfo PersonalInfo `yaml:",inline" json:",inline"`

	MonthlySalaries []SalaryBase `yaml:"monthly_salaries" json:"monthly_salaries"`
}

// MonthlyTaxes 返回的对象
type MonthlyTaxes struct {
	Taxes []*MonthlyTax `yaml:"taxes" json:"taxes"`
}

// MonthlyTax 月薪对象
type MonthlyTax struct {
	Month int `yaml:"month" json:"month"`

	SalaryBase `yaml:",inline" json:",inline"`

	InsurancesResult       *CalcInsurancesAmount `yaml:"insurances_result" json:"insurances_result"`
	AccumulationFundResult *CalcAccumulationFund `yaml:"accumulation_fund_result" json:"accumulation_fund_result"`

	Insurances       float64 `yaml:"insurances" json:"insurances"`
	AccumulationFund float64 `yaml:"accumulation_fund" json:"accumulation_fund"`

	Taxation        float64 `yaml:"taxation" json:"taxation"`
	RestSalary      float64 `yaml:"rest_salary" json:"rest_salary"`
	HistoryTaxation float64 `yaml:"history_taxation" json:"history_taxation"`
	HistorySalary   float64 `yaml:"history_salary" json:"history_salary"`
}

// Calc 计算月薪剩余以及个税情况
func (p *TaxesHandler) Calc(salaries *Salaries) (t *MonthlyTaxes, err error) {
	taxes := &MonthlyTaxes{}
	info := &PersonalInfo{}
	lastMonth := 0
	for i, s := range salaries.MonthlySalaries {
		lastMonth = i + 1
		iMonthTax := &MonthlyTax{
			Month:      lastMonth,
			SalaryBase: s,
		}

		info.SalaryBase = s

		iMonthTax.AccumulationFundResult, err = p.AccumulationFundHandler.Calc(info)
		if err != nil {
			return nil, err
		}

		iMonthTax.InsurancesResult, err = p.InsurancesHandler.Calc(info)
		if err != nil {
			return nil, err
		}

		iMonthTax.Insurances = iMonthTax.InsurancesResult.Private.EndowmentAmount +
			iMonthTax.InsurancesResult.Private.MedicalAmount +
			iMonthTax.InsurancesResult.Private.UnemploymentAmount +
			iMonthTax.InsurancesResult.Private.EmploymentInjuryAmount +
			iMonthTax.InsurancesResult.Private.BirthAmount +
			iMonthTax.InsurancesResult.Private.SeriousMedicalAmount

		iMonthTax.AccumulationFund = iMonthTax.AccumulationFundResult.PrivateFund

		p.getMonthTax(iMonthTax)

		taxes.Taxes = append(taxes.Taxes, iMonthTax)
	}

	if salaries.For {
		for i := lastMonth; i < 12; i++ {
			monthlyTax := taxes.Taxes[lastMonth-1]

			monthlyTaxNext := *monthlyTax
			monthlyTaxNext.Month = i + 1
			p.getMonthTax(&monthlyTaxNext)
			taxes.Taxes = append(taxes.Taxes, &monthlyTaxNext)
		}
	}

	return taxes, nil
}

func (p *TaxesHandler) getMonthTax(monthlyTax *MonthlyTax) {

	monthlyTax.RestSalary = Decimal2(monthlyTax.Salary + monthlyTax.SubsidyAmount -
		monthlyTax.Insurances - monthlyTax.AccumulationFund)
	taxSalary := monthlyTax.RestSalary - monthlyTax.Threshold - monthlyTax.DeductibleAmount
	if taxSalary > 0 {
		p.totalTaxSalaries += taxSalary
	}
	p.totalSalaries += monthlyTax.RestSalary

	// 小于起征点，那么税收为0
	if monthlyTax.RestSalary-monthlyTax.DeductibleAmount <= monthlyTax.Threshold {
		return
	}

	for _, taxRate := range p.YearTaxRates {
		if p.totalTaxSalaries <= taxRate.SalaryMin ||
			(p.totalTaxSalaries > taxRate.SalaryMax && taxRate.SalaryMax != 0) {
			continue
		}
		tax := Decimal2((p.totalTaxSalaries*taxRate.Rate)/100.0 - taxRate.DeductedAmount - p.totalTaxation)
		p.totalTaxation += tax

		monthlyTax.Taxation = tax
		monthlyTax.RestSalary = Decimal2(monthlyTax.RestSalary - tax)
		monthlyTax.HistorySalary = Decimal2(p.totalSalaries)
		monthlyTax.HistoryTaxation = Decimal2(p.totalTaxation)
		return
	}
}

const (
	printTaxInfor = "%2d月, 收入: %10.2f, 补贴: %10.2f, 社保缴纳: %4.2f, 公积金缴纳: %4.2f, 个税缴纳: %10.2f, 剩余工资: %10.2f"
)

// Print 打印信息
func (p *MonthlyTaxes) Print() {
	for _, t := range p.Taxes {
		fmt.Println(fmt.Sprintf(printTaxInfor, t.Month, t.Salary, t.SubsidyAmount,
			t.Insurances, t.AccumulationFund, t.Taxation, t.RestSalary))
	}
}
