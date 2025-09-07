package models

type Match struct {
	ID            int    `json:"match_num"`
	IsRedTable    bool   `json:"is_red_table"`
	AllianceA     int    `json:"alliance_a"`
	AllianceAName string `json:"alliance_a_name"`
	AllianceB     int    `json:"alliance_b"`
	AllianceBName string `json:"alliance_b_name"`
}
