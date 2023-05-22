package dao

import "time"

type TriggersModel struct {
	Modon time.Time `json:"modon"`
	// tigger id: templateId:randomStr(3)
	Triggers map[string]string `json:"triggers"`
}

type TemplatesModel struct {
	Modon time.Time `json:"modon"`
	// tigger id: templateId:randomStr(3)
	Templates map[string]string `json:"templates"`
}
