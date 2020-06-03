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

// Insurances 社保对象
type Insurances struct {
	InsurancesBase `yaml:"insurances" json:"insurances"`
}

// InsurancesAmount 社保金额
type InsurancesAmount struct {
	// 基本养老
	EndowmentAmount    float64 `yaml:"endowment_amount" json:"endowment_amount"`
	MaxEndowmentAmount float64 `yaml:"max_endowment_amount" json:"max_endowment_amount"`
	MinEndowmentAmount float64 `yaml:"min_endowment_amount" json:"min_endowment_amount"`
	// 基本医保
	MedicalAmount    float64 `yaml:"medical_amount" json:"medical_amount"`
	MaxMedicalAmount float64 `yaml:"max_medical_amount" json:"max_medical_amount"`
	MinMedicalAmount float64 `yaml:"min_medical_amount" json:"min_medical_amount"`
	// 失业
	UnemploymentAmount    float64 `yaml:"unemployment_amount" json:"unemployment_amount"`
	MaxUnemploymentAmount float64 `yaml:"max_unemployment_amount" json:"max_unemployment_amount"`
	MinUnemploymentAmount float64 `yaml:"min_unemployment_amount" json:"min_unemployment_amount"`
	// 工伤
	EmploymentInjuryAmount    float64 `yaml:"employment_injury_amount" json:"employment_injury_amount"`
	MaxEmploymentInjuryAmount float64 `yaml:"max_employment_injury_amount" json:"max_employment_injury_amount"`
	MinEmploymentInjuryAmount float64 `yaml:"min_employment_injury_amount" json:"min_employment_injury_amount"`
	// 生育
	BirthAmount    float64 `yaml:"birth_amount" json:"birth_amount"`
	MaxBirthAmount float64 `yaml:"max_birth_amount" json:"max_birth_amount"`
	MinBirthAmount float64 `yaml:"min_birth_amount" json:"min_birth_amount"`
	// 大病医保
	SeriousMedicalAmount    float64 `yaml:"serious_medical_amount" json:"serious_medical_amount"`
	MaxSeriousMedicalAmount float64 `yaml:"max_serious_medical_amount" json:"max_serious_medical_amount"`
	MinSeriousMedicalAmount float64 `yaml:"min_serious_medical_amount" json:"min_serious_medical_amount"`
}

// CalcInsurancesAmount 结果对象
type CalcInsurancesAmount struct {
	InsurancesBase `yaml:"insurances" json:"insurances"`
	PersonalInfo   `yaml:"persion_info" json:"persion_info"`

	Company InsurancesAmount `yaml:"company_insurances" json:"company_insurances"`
	Private InsurancesAmount `yaml:"private_insurances" json:"private_insurances"`

	CompanyTotalAmount float64 `yaml:"company_total_amount" json:"company_total_amount"`
	PrivateTotalAmount float64 `yaml:"private_total_amount" json:"private_total_amount"`
}

// NewInsurances 生成社保对象
func NewInsurances(file string) (*Insurances, error) {
	i := &Insurances{}
	err := config.NewSuffixReader().Read(file, i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

const (
	printInsuranceInfor = "%s, 最低基数: %.2f, 最高基数: %.2f, 单位承担比例: %.2f%%, 个人承担比例: %.2f%%, 单位最低金额: %.2f, 单位最高金额: %.2f, 个人最低金额: %.2f, 个人最高金额: %.2f. \n\t  实际基数: %.2f, 单位缴纳: %.2f, 个人缴纳: %.2f"
)

// Print 打印信息
func (p *CalcInsurancesAmount) Print() {
	switch p.PersonalInfo.Endowment {
	case EndowmentOffice:
		fmt.Println(fmt.Sprintf(printInsuranceInfor, "基本养老",
			p.InsurancesBase.OfficeEndowment.MinBase, p.InsurancesBase.OfficeEndowment.MaxBase,
			p.InsurancesBase.OfficeEndowment.CompanyRate, p.InsurancesBase.OfficeEndowment.PrivateRate,
			p.Company.MinEndowmentAmount, p.Company.MaxEndowmentAmount,
			p.Private.MinEndowmentAmount, p.Private.MaxEndowmentAmount,
			p.PersonalInfo.EndowmentBase, p.Company.EndowmentAmount, p.Private.EndowmentAmount,
		))
	default:
		fmt.Println(fmt.Sprintf(printInsuranceInfor, "基本养老",
			p.InsurancesBase.WorkersEndowment.MinBase, p.InsurancesBase.WorkersEndowment.MaxBase,
			p.InsurancesBase.WorkersEndowment.CompanyRate, p.InsurancesBase.WorkersEndowment.PrivateRate,
			p.Company.MinEndowmentAmount, p.Company.MaxEndowmentAmount,
			p.Private.MinEndowmentAmount, p.Private.MaxEndowmentAmount,
			p.PersonalInfo.EndowmentBase, p.Company.EndowmentAmount, p.Private.EndowmentAmount,
		))
	}

	fmt.Println(fmt.Sprintf(printInsuranceInfor, "基本医疗",
		p.InsurancesBase.Medical.MinBase, p.InsurancesBase.Medical.MaxBase,
		p.InsurancesBase.Medical.CompanyRate, p.InsurancesBase.Medical.PrivateRate,
		p.Company.MinMedicalAmount, p.Company.MaxMedicalAmount,
		p.Private.MinMedicalAmount, p.Private.MaxMedicalAmount,
		p.PersonalInfo.MedicalBase, p.Company.MedicalAmount, p.Private.MedicalAmount,
	))

	switch p.PersonalInfo.Residence {
	case ResidenceNonAgricultural:
		fmt.Println(fmt.Sprintf(printInsuranceInfor, "失业保险",
			p.InsurancesBase.NonAgriculturalUnemployment.MinBase, p.InsurancesBase.NonAgriculturalUnemployment.MaxBase,
			p.InsurancesBase.NonAgriculturalUnemployment.CompanyRate, p.InsurancesBase.NonAgriculturalUnemployment.PrivateRate,
			p.Company.MinUnemploymentAmount, p.Company.MaxUnemploymentAmount,
			p.Private.MinUnemploymentAmount, p.Private.MaxUnemploymentAmount,
			p.PersonalInfo.UnemploymentBase, p.Company.UnemploymentAmount, p.Private.UnemploymentAmount,
		))
	default:
		fmt.Println(fmt.Sprintf(printInsuranceInfor, "失业保险",
			p.InsurancesBase.AgriculturalUnemployment.MinBase, p.InsurancesBase.AgriculturalUnemployment.MaxBase,
			p.InsurancesBase.AgriculturalUnemployment.CompanyRate, p.InsurancesBase.AgriculturalUnemployment.PrivateRate,
			p.Company.MinUnemploymentAmount, p.Company.MaxUnemploymentAmount,
			p.Private.MinUnemploymentAmount, p.Private.MaxUnemploymentAmount,
			p.PersonalInfo.UnemploymentBase, p.Company.UnemploymentAmount, p.Private.UnemploymentAmount,
		))
	}
	fmt.Println(fmt.Sprintf(printInsuranceInfor, "工伤保险",
		p.InsurancesBase.EmploymentInjury.MinBase, p.InsurancesBase.EmploymentInjury.MaxBase,
		p.InsurancesBase.EmploymentInjury.CompanyRate, p.InsurancesBase.EmploymentInjury.PrivateRate,
		p.Company.MinEmploymentInjuryAmount, p.Company.MaxEmploymentInjuryAmount,
		p.Private.MinEmploymentInjuryAmount, p.Private.MaxEmploymentInjuryAmount,
		p.PersonalInfo.EmploymentInjuryBase, p.Company.EmploymentInjuryAmount, p.Private.EmploymentInjuryAmount,
	))

	fmt.Println(fmt.Sprintf(printInsuranceInfor, "生育保险",
		p.InsurancesBase.Birth.MinBase, p.InsurancesBase.Birth.MaxBase,
		p.InsurancesBase.Birth.CompanyRate, p.InsurancesBase.Birth.PrivateRate,
		p.Company.MinBirthAmount, p.Company.MaxBirthAmount,
		p.Private.MinBirthAmount, p.Private.MaxBirthAmount,
		p.PersonalInfo.BirthBase, p.Company.BirthAmount, p.Private.BirthAmount,
	))

	fmt.Println(fmt.Sprintf(printInsuranceInfor, "大病医疗",
		p.InsurancesBase.SeriousMedical.MinBase, p.InsurancesBase.SeriousMedical.MaxBase,
		p.InsurancesBase.SeriousMedical.CompanyRate, p.InsurancesBase.SeriousMedical.PrivateRate,
		p.Company.MinSeriousMedicalAmount, p.Company.MaxSeriousMedicalAmount,
		p.Private.MinSeriousMedicalAmount, p.Private.MaxSeriousMedicalAmount,
		p.PersonalInfo.SeriousMedicalBase, p.Company.SeriousMedicalAmount, p.Private.SeriousMedicalAmount,
	))

	fmt.Println(fmt.Sprintf("\t单位总承担: %0.2f, 个人总承担: %0.2f", p.CompanyTotalAmount, p.PrivateTotalAmount))
}

// Calc 计算
func (p *Insurances) Calc(info *PersonalInfo) (*CalcInsurancesAmount, error) {
	calc := &CalcInsurancesAmount{
		InsurancesBase: p.InsurancesBase,
		PersonalInfo:   *info,
	}
	switch info.Endowment {
	case EndowmentOffice:
		calc.Company.EndowmentAmount = Decimal2(info.EndowmentBase * p.OfficeEndowment.CompanyRate / 100.0)
		calc.Company.MinEndowmentAmount = Decimal2(p.OfficeEndowment.MinBase * p.OfficeEndowment.CompanyRate / 100.0)
		calc.Company.MaxEndowmentAmount = Decimal2(p.OfficeEndowment.MaxBase * p.OfficeEndowment.CompanyRate / 100.0)

		calc.Private.EndowmentAmount = Decimal2(info.EndowmentBase * p.OfficeEndowment.PrivateRate / 100.0)
		calc.Private.MinEndowmentAmount = Decimal2(p.OfficeEndowment.MinBase * p.OfficeEndowment.PrivateRate / 100.0)
		calc.Private.MaxEndowmentAmount = Decimal2(p.OfficeEndowment.MaxBase * p.OfficeEndowment.PrivateRate / 100.0)
	default:
		// WorkersEndowment
		calc.Company.EndowmentAmount = Decimal2(info.EndowmentBase * p.WorkersEndowment.CompanyRate / 100.0)
		calc.Company.MinEndowmentAmount = Decimal2(p.WorkersEndowment.MinBase * p.WorkersEndowment.CompanyRate / 100.0)
		calc.Company.MaxEndowmentAmount = Decimal2(p.WorkersEndowment.MaxBase * p.WorkersEndowment.CompanyRate / 100.0)

		calc.Private.EndowmentAmount = Decimal2(info.EndowmentBase * p.WorkersEndowment.PrivateRate / 100.0)
		calc.Private.MinEndowmentAmount = Decimal2(p.WorkersEndowment.MinBase * p.WorkersEndowment.PrivateRate / 100.0)
		calc.Private.MaxEndowmentAmount = Decimal2(p.WorkersEndowment.MaxBase * p.WorkersEndowment.PrivateRate / 100.0)
	}

	calc.CompanyTotalAmount += calc.Company.EndowmentAmount
	calc.PrivateTotalAmount += calc.Private.EndowmentAmount

	switch info.Residence {
	case ResidenceNonAgricultural:
		// NonAgriculturalUnemployment
		calc.Company.UnemploymentAmount =
			Decimal2(info.UnemploymentBase * p.NonAgriculturalUnemployment.CompanyRate / 100.0)
		calc.Company.MinUnemploymentAmount =
			Decimal2(p.NonAgriculturalUnemployment.MinBase * p.NonAgriculturalUnemployment.CompanyRate / 100.0)
		calc.Company.MaxUnemploymentAmount =
			Decimal2(p.NonAgriculturalUnemployment.MaxBase * p.NonAgriculturalUnemployment.CompanyRate / 100.0)

		calc.Private.UnemploymentAmount =
			Decimal2(info.UnemploymentBase * p.NonAgriculturalUnemployment.PrivateRate / 100.0)
		calc.Private.MinUnemploymentAmount =
			Decimal2(p.NonAgriculturalUnemployment.MinBase * p.NonAgriculturalUnemployment.PrivateRate / 100.0)
		calc.Private.MaxUnemploymentAmount =
			Decimal2(p.NonAgriculturalUnemployment.MaxBase * p.NonAgriculturalUnemployment.PrivateRate / 100.0)
	default:
		// AgriculturalUnemployment
		calc.Company.UnemploymentAmount =
			Decimal2(info.UnemploymentBase * p.AgriculturalUnemployment.CompanyRate / 100.0)
		calc.Company.MinUnemploymentAmount =
			Decimal2(p.AgriculturalUnemployment.MinBase * p.AgriculturalUnemployment.CompanyRate / 100.0)
		calc.Company.MaxUnemploymentAmount =
			Decimal2(p.AgriculturalUnemployment.MaxBase * p.AgriculturalUnemployment.CompanyRate / 100.0)

		calc.Private.UnemploymentAmount =
			Decimal2(info.UnemploymentBase * p.AgriculturalUnemployment.PrivateRate / 100.0)
		calc.Private.MinUnemploymentAmount =
			Decimal2(p.AgriculturalUnemployment.MinBase * p.AgriculturalUnemployment.PrivateRate / 100.0)
		calc.Private.MaxUnemploymentAmount =
			Decimal2(p.AgriculturalUnemployment.MaxBase * p.AgriculturalUnemployment.PrivateRate / 100.0)
	}

	calc.CompanyTotalAmount += calc.Company.UnemploymentAmount
	calc.PrivateTotalAmount += calc.Private.UnemploymentAmount

	calc.Company.EmploymentInjuryAmount =
		Decimal2(info.EmploymentInjuryBase * p.EmploymentInjury.CompanyRate / 100.0)
	calc.Company.MinEmploymentInjuryAmount =
		Decimal2(p.EmploymentInjury.MinBase * p.EmploymentInjury.CompanyRate / 100.0)
	calc.Company.MaxEmploymentInjuryAmount =
		Decimal2(p.EmploymentInjury.MaxBase * p.EmploymentInjury.CompanyRate / 100.0)
	calc.Private.EmploymentInjuryAmount =
		Decimal2(info.EmploymentInjuryBase * p.EmploymentInjury.PrivateRate / 100.0)
	calc.Private.MinEmploymentInjuryAmount =
		Decimal2(p.EmploymentInjury.MinBase * p.EmploymentInjury.PrivateRate / 100.0)
	calc.Private.MaxEmploymentInjuryAmount =
		Decimal2(p.EmploymentInjury.MaxBase * p.EmploymentInjury.PrivateRate / 100.0)

	calc.CompanyTotalAmount += calc.Company.EmploymentInjuryAmount
	calc.PrivateTotalAmount += calc.Private.EmploymentInjuryAmount

	calc.Company.BirthAmount = Decimal2(info.BirthBase * p.Birth.CompanyRate / 100.0)
	calc.Company.MinBirthAmount = Decimal2(p.Birth.MinBase * p.Birth.CompanyRate / 100.0)
	calc.Company.MaxBirthAmount = Decimal2(p.Birth.MaxBase * p.Birth.CompanyRate / 100.0)
	calc.Private.BirthAmount = Decimal2(info.BirthBase * p.Birth.PrivateRate / 100.0)
	calc.Private.MinBirthAmount = Decimal2(p.Birth.MinBase * p.Birth.PrivateRate / 100.0)
	calc.Private.MaxBirthAmount = Decimal2(p.Birth.MaxBase * p.Birth.PrivateRate / 100.0)

	calc.CompanyTotalAmount += calc.Company.BirthAmount
	calc.PrivateTotalAmount += calc.Private.BirthAmount

	calc.Company.MedicalAmount = Decimal2(info.MedicalBase * p.Medical.CompanyRate / 100.0)
	calc.Company.MinMedicalAmount = Decimal2(p.Medical.MinBase * p.Medical.CompanyRate / 100.0)
	calc.Company.MaxMedicalAmount = Decimal2(p.Medical.MaxBase * p.Medical.CompanyRate / 100.0)
	calc.Private.MedicalAmount = Decimal2(info.MedicalBase * p.Medical.PrivateRate / 100.0)
	calc.Private.MinMedicalAmount = Decimal2(p.Medical.MinBase * p.Medical.PrivateRate / 100.0)
	calc.Private.MaxMedicalAmount = Decimal2(p.Medical.MaxBase * p.Medical.PrivateRate / 100.0)

	calc.CompanyTotalAmount += calc.Company.MedicalAmount
	calc.PrivateTotalAmount += calc.Private.MedicalAmount

	calc.Company.SeriousMedicalAmount = Decimal2(info.SeriousMedicalBase * p.SeriousMedical.CompanyRate / 100.0)
	calc.Company.MinSeriousMedicalAmount = Decimal2(p.SeriousMedical.MinBase * p.SeriousMedical.CompanyRate / 100.0)
	calc.Company.MaxSeriousMedicalAmount = Decimal2(p.SeriousMedical.MaxBase * p.SeriousMedical.CompanyRate / 100.0)
	calc.Private.SeriousMedicalAmount =
		Decimal2(info.SeriousMedicalBase*p.SeriousMedical.PrivateRate/100.0 + p.SeriousMedical.ExtraPayment)
	calc.Private.MinSeriousMedicalAmount =
		Decimal2(p.SeriousMedical.MinBase*p.SeriousMedical.PrivateRate/100.0 + p.SeriousMedical.ExtraPayment)
	calc.Private.MaxSeriousMedicalAmount =
		Decimal2(p.SeriousMedical.MaxBase*p.SeriousMedical.PrivateRate/100.0 + p.SeriousMedical.ExtraPayment)

	calc.CompanyTotalAmount += calc.Company.SeriousMedicalAmount
	calc.PrivateTotalAmount += calc.Private.SeriousMedicalAmount

	return calc, nil
}
