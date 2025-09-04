package models

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WebSocketResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WebSocketMatchPayload struct {
	MatchID      int      `json:"match_id"`
	IsRedTable   bool     `json:"is_red_table"`
	AllianceA  []string `json:"alliance_a"`
	AllianceB []string `json:"alliance_b"`
}
