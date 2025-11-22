package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jake-schuler/fll-event-screen/config"
	"github.com/jake-schuler/fll-event-screen/models"
)

type scheduleRound int

const (
	practice scheduleRound = iota
	round1
	round2
	round3
)

var ScheduleRound = map[scheduleRound]string{
	practice: "Practice",
	round1:   "Round 1",
	round2:   "Round 2",
	round3:   "Round 3",
}

func ReadSchedule(scheduleRound scheduleRound) []models.Match {
	// Read whole schedule area (header at row 1)
	resp, err := config.Srv.Spreadsheets.Values.Get(config.SheetID, "Schedule!A2:E").Do()
	if err != nil {
		fmt.Printf("Error reading schedule: %v\n", err)
		return []models.Match{}
	}

	// Map keyed by round|time|color (Red/Blue). Each key represents one match (A and B sides)
	matchMap := make(map[string]*models.Match)
	keys := make([]string, 0)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return []models.Match{}
	}

	for _, row := range resp.Values {
		// Expecting: Team #, Team Name, Round, Match Start, Table
		// tolerate shorter rows
		if len(row) == 0 {
			continue
		}

		// Team number may be empty (wild card) or string/number
		var teamNum int
		if len(row) > 0 {
			switch v := row[0].(type) {
			case string:
				v = strings.TrimSpace(v)
				if v == "" {
					teamNum = 0
				} else {
					if n, err := strconv.Atoi(v); err == nil {
						teamNum = n
					} else {
						// not a number -> treat as 0
						teamNum = 0
					}
				}
			case float64:
				teamNum = int(v)
			default:
				teamNum = 0
			}
		}

		// Team name
		var teamName string
		if len(row) > 1 {
			if s, ok := row[1].(string); ok {
				teamName = strings.TrimSpace(s)
			}
		}

		// Round (e.g., "Practice", "Round 1")
		var roundStr string
		if len(row) > 2 {
			if s, ok := row[2].(string); ok {
				roundStr = strings.TrimSpace(s)
			}
		}

		// Match start time string (e.g., "9:30 AM")
		var matchStart string
		if len(row) > 3 {
			if s, ok := row[3].(string); ok {
				matchStart = strings.TrimSpace(s)
			}
		}

		// Table (e.g., "Red A", "Blue B")
		var table string
		if len(row) > 4 {
			if s, ok := row[4].(string); ok {
				table = strings.TrimSpace(s)
			}
		}

		if table == "" {
			// skip rows without table assignment
			continue
		}

		color := ""
		if strings.Contains(strings.ToLower(table), "red") {
			color = "Red"
		} else if strings.Contains(strings.ToLower(table), "blue") {
			color = "Blue"
		} else {
			// unknown color, skip
			continue
		}

		side := ""
		if strings.HasSuffix(table, "A") {
			side = "A"
		} else if strings.HasSuffix(table, "B") {
			side = "B"
		} else {
			// unknown side -> skip
			continue
		}

		// key groups by round + start time + color
		key := fmt.Sprintf("%s|%s|%s", roundStr, matchStart, color)

		m, exists := matchMap[key]
		if !exists {
			m = &models.Match{
				ID:            0,
				IsRedTable:    color == "Red",
				AllianceA:     0,
				AllianceAName: "",
				AllianceB:     0,
				AllianceBName: "",
			}
			matchMap[key] = m
			keys = append(keys, key)
		}

		if side == "A" {
			m.AllianceA = teamNum
			m.AllianceAName = teamName
		} else if side == "B" {
			m.AllianceB = teamNum
			m.AllianceBName = teamName
		}
	}

	// Preserve sheet order: keys were appended in the order groups were first encountered
	matches := make([]models.Match, 0, len(keys))
	id := 1
	for _, key := range keys {
		m := matchMap[key]
		m.ID = id
		matches = append(matches, *m)
		id++
	}

	return matches
}
