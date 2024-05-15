package entity

type Fact struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMoID     int    `json:"indicator_to_mo_id"`
	IndicatorToMoFactID int    `json:"indicator_to_mo_fact_id"`
	Value               int    `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              int    `json:"is_plan"`
	AuthUserID          int    `json:"auth_user_id"`
	Comment             string `json:"comment"`
}
