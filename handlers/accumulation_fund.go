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

// AccumulationFund 公积金对象
type AccumulationFund struct {
	AccumulationFundBase `yaml:"accumulation_fund" json:"accumulation_fund"`
}

// CalcAccumulationFund 结果对象
type CalcAccumulationFund struct {
	AccumulationFundBase `yaml:"accumulation_fund" json:"accumulation_fund"`

	Salary float64 `yaml:"salary" json:"salary"`
	Base   float64 `yaml:"base" json:"base"`
	Rate   float64 `yaml:"rate" json:"rate"`

	CompanyFund    float64 `yaml:"company_fund" json:"company_fund"`
	MinCompanyFund float64 `yaml:"min_company_fund" json:"min_company_fund"`
	MaxCompanyFund float64 `yaml:"max_company_fund" json:"max_company_fund"`
	PrivateFund    float64 `yaml:"private_fund" json:"private_fund"`
	MinPrivateFund float64 `yaml:"min_private_fund" json:"min_private_fund"`
	MaxPrivateFund float64 `yaml:"max_private_fund" json:"max_private_fund"`
}

// NewAccumulationFund 生成公积金对象
func NewAccumulationFund(file string) (*AccumulationFund, error) {
	a := &AccumulationFund{}
	err := config.NewSuffixReader().Read(file, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

const (
	printFundInfor = "%s, 最低基数: %0.f, 最高基数: %0.f, 最低比例: %0.2f%%, 最高比例: %0.2f%%, 单位最低金额: %0.f, 单位最高金额: %0.f, 个人最低金额: %0.f, 个人最高金额: %0.f. \n\t  实际基数: %0.f, 缴纳比例: %0.2f%%, 单位缴纳: %0.f, 个人缴纳: %0.f"
)

// Print 打印信息
func (p *CalcAccumulationFund) Print() {
	fmt.Println(fmt.Sprintf(printFundInfor, "公积金",
		p.AccumulationFundBase.MinBase, p.AccumulationFundBase.MaxBase,
		p.AccumulationFundBase.MinRate, p.AccumulationFundBase.MaxRate,
		p.MinCompanyFund, p.MaxCompanyFund,
		p.MinPrivateFund, p.MaxPrivateFund,
		p.Base, p.Rate, p.CompanyFund, p.PrivateFund,
	))
}

// Calc 计算
func (p *AccumulationFund) Calc(info *PersonalInfo) (*CalcAccumulationFund, error) {
	result := &CalcAccumulationFund{
		AccumulationFundBase: p.AccumulationFundBase,

		Salary: info.Salary,
		Rate:   info.AccumulationFundRate,
	}

	if result.Salary > p.AccumulationFundBase.MaxBase {
		result.Base = p.AccumulationFundBase.MaxBase
	} else if result.Salary < p.AccumulationFundBase.MinBase {
		return result, fmt.Errorf("小于最小基数: %0.f", p.AccumulationFundBase.MinBase)
	} else {
		result.Base = info.Salary
	}

	if result.Rate > p.AccumulationFundBase.MaxRate || result.Rate < p.AccumulationFundBase.MinRate {
		return nil, fmt.Errorf("比例需在 %0.f 和 %0.f 之间",
			p.AccumulationFundBase.MinRate, p.AccumulationFundBase.MaxRate)
	}

	result.CompanyFund = Decimal(result.Base*result.Rate/100.0, 0)
	result.PrivateFund = result.CompanyFund

	result.MaxCompanyFund = Decimal(result.MaxBase*result.MaxRate/100.0, 0)
	result.MinCompanyFund = Decimal(result.MinBase*result.MinRate/100.0, 0)
	result.MaxPrivateFund = result.MaxCompanyFund
	result.MinPrivateFund = result.MinCompanyFund

	return result, nil
}
