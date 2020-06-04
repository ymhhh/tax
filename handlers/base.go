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
	"strconv"
)

// InsurancesBase 社保信息
type InsurancesBase struct {
	// 职工基本养老
	WorkersEndowment Base `yaml:"workers_endowment" json:"workers_endowment"`
	// 机关基本养老
	OfficeEndowment Base `yaml:"office_endowment" json:"office_endowment"`
	// 基本医保
	Medical Base `yaml:"medical" json:"medical"`
	// 非农失业
	NonAgriculturalUnemployment Base `yaml:"non_agricultural_unemployment" json:"non_agricultural_unemployment"`
	// 农业失业
	AgriculturalUnemployment Base `yaml:"agricultural_unemployment" json:"agricultural_unemployment"`
	// 工伤
	EmploymentInjury Base `yaml:"employment_injury" json:"employment_injury"`
	// 生育
	Birth Base `yaml:"birth" json:"birth"`
	// 大病医保
	SeriousMedical Base `yaml:"serious_medical" json:"serious_medical"`
}

// Base 基础信息
type Base struct {
	MinBase      float64 `yaml:"min_base" json:"min_base"`
	MaxBase      float64 `yaml:"max_base" json:"max_base"`
	PrivateRate  float64 `yaml:"private_rate" json:"private_rate"`
	CompanyRate  float64 `yaml:"company_rate" json:"company_rate"`
	ExtraPayment float64 `yaml:"extra_payment" json:"extra_payment"`
}

// AccumulationFundBase 公积金基数
type AccumulationFundBase struct {
	MinBase float64 `yaml:"min_base" json:"min_base"`
	MaxBase float64 `yaml:"max_base" json:"max_base"`
	MinRate float64 `yaml:"min_rate" json:"min_rate"`
	MaxRate float64 `yaml:"max_rate" json:"max_rate"`
}

// ResidenceType 定义户口类型
type ResidenceType int

// Residence 户口类型
const (
	// 非农户口
	ResidenceNonAgricultural ResidenceType = iota
	// 农业户口
	ResidenceAgricultural
)

// EndowmentType 定义养老类型
type EndowmentType int

// 养老类型
const (
	// 普通职工
	EndowmentWorkers EndowmentType = iota
	// 机关单位
	EndowmentOffice
)

// PersonalInfo 个人基本信息
type PersonalInfo struct {
	Residence ResidenceType `yaml:"residence" json:"residence"`
	Endowment EndowmentType `yaml:"endowment" json:"endowment"`

	SalaryBase `yaml:",inline" json:",inline"`
}

// SalaryBase 基本薪水信息
type SalaryBase struct {
	Threshold        float64 `yaml:"threshold" json:"threshold"`                 // 基数
	Salary           float64 `yaml:"salary" json:"salary"`                       // 薪水
	SubsidyAmount    float64 `yaml:"subsidy_amount" json:"subsidy_amount"`       // 补贴
	DeductibleAmount float64 `yaml:"deductible_amount" json:"deductible_amount"` // 抵扣金额

	AccumulationFundRate float64 `yaml:"accumulation_fund_rate" json:"accumulation_fund_rate"`

	EndowmentBase        float64 `yaml:"endowment_base" json:"endowment_base"`
	MedicalBase          float64 `yaml:"medical_base" json:"medical_base"`
	UnemploymentBase     float64 `yaml:"unemployment_base" json:"unemployment_base"`
	EmploymentInjuryBase float64 `yaml:"employment_injury_base" json:"employment_injury_base"`
	BirthBase            float64 `yaml:"birth_base" json:"birth_base"`
	SeriousMedicalBase   float64 `yaml:"serious_medical_base" json:"serious_medical_base"`
}

// Decimal 处理浮点数精度
func Decimal(value float64, pos int) float64 {
	format := fmt.Sprintf("%%.%df", pos)
	value, _ = strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return value
}

// Decimal2 2位浮点精度
func Decimal2(value float64) float64 {
	return Decimal(value, 2)
}

// YearTaxBase 个税年情况
type YearTaxBase struct {
	YearTaxRates []YearTaxRate `yaml:"year_tax_rates" json:"year_tax_rates"`
}

// YearTaxRate 年配置
type YearTaxRate struct {
	SalaryMin      float64 `yaml:"salary_min" json:"salary_min"`
	SalaryMax      float64 `yaml:"salary_max" json:"salary_max"`
	Rate           float64 `yaml:"rate" json:"rate"`
	DeductedAmount float64 `yaml:"deducted_amount" json:"deducted_amount"`
}
