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
	"github.com/go-trellis/config"
)

type Taxes struct {
	AccumulationFund `yaml:",inline" json:",inline"`
	Insurances       `yaml:",inline" json:",inline"`

	YearTaxBase `yaml:",inline" json:",inline"`

	totalSalaries    float64 `yaml:"-" json:"-"`
	totalTaxSalaries float64 `yaml:"-" json:"-"`
	totalTaxation    float64 `yaml:"-" json:"-"`
}

func NewTaxes(file string) (*Taxes, error) {
	t := &Taxes{}
	if err := config.NewSuffixReader().Read(file, t); err != nil {
		return nil, err
	}
	return t, nil
}

type Salaries struct {
	For bool `yaml:"for" json:"for"`

	PersonalInfo PersonalInfo `yaml:",inline" json:",inline"`

	MonthlySalaries []SalaryBase `yaml:"monthly_salaries" json:"monthly_salaries"`
}

type MonthlyTaxes struct {
	Taxes []*MonthlyTax `yaml:"taxes" json:"taxes"`
}

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
func (p *Taxes) Calc(salaries *Salaries) (t *MonthlyTaxes, err error) {
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

		iMonthTax.AccumulationFundResult, err = p.AccumulationFund.Calc(info)
		if err != nil {
			return nil, err
		}

		iMonthTax.InsurancesResult, err = p.Insurances.Calc(info)
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

		taxes.Taxes = append(taxes.Taxes, iMonthTax)
	}

	if salaries.For {
		for i := lastMonth + 1; i < 13; i++ {
			monthlyTax := taxes.Taxes[lastMonth-1]

			monthlyTaxNext := *monthlyTax
			monthlyTaxNext.Month = i
			p.getMonthTax(i, &monthlyTaxNext)
			taxes.Taxes = append(taxes.Taxes, &monthlyTaxNext)
		}
	}

	return taxes, nil
}

func (p *Taxes) getMonthTax(lastMonth int, monthlyTax *MonthlyTax) {

	monthlyTax.RestSalary = monthlyTax.Salary + monthlyTax.SubsidyAmount -
		monthlyTax.Insurances - monthlyTax.AccumulationFund
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
		tax := (p.totalTaxSalaries*taxRate.Rate)/100.0 - taxRate.DeductedAmount - p.totalTaxation
		p.totalTaxation += tax

		monthlyTax.Taxation = tax
		monthlyTax.RestSalary = monthlyTax.RestSalary - tax
		monthlyTax.HistorySalary = p.totalSalaries
		monthlyTax.HistoryTaxation = p.totalTaxation
		return
	}
}
