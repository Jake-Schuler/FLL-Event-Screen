package models

type Match struct {
	ID int `json:"match_num"`
	IsRedTable bool `json:"is_red_table"`
	AllianceA []string `json:"alliance_a"`
	AllianceB []string `json:"alliance_b"`
}