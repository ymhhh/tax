package handlers

type InsurancesBase struct {
	Endowment        Base `yaml:"endowment" json:"endowment"`                 // 基本养老
	Medical          Base `yaml:"medical" json:"medical"`                     // 基本医保
	Unemployment     Base `yaml:"unemployment" json:"unemployment"`           // 失业
	EmploymentInjury Base `yaml:"employment_injury" json:"employment_injury"` // 工伤
	Birth            Base `yaml:"birth" json:"birth"`                         // 生育
	SeriousMedical   Base `yaml:"serious_medical" json:"serious_medical"`     // 大病医保
}

type Base struct {
	Base        float64 `yaml:"base" json:"base"`
	Minbase     float64 `yaml:"min_base" json:"min_base"`
	Maxbase     float64 `yaml:"max_base" json:"max_base"`
	PrivateRate float64 `yaml:"private_rate" json:"private_rate"`
	CompanyRate float64 `yaml:"company_rate" json:"company_rate"`
	MinPayment  float64 `yaml:"min_payment" json:"min_payment"`
	MaxPayment  float64 `yaml:"max_payment" json:"max_payment"`
}
