package services

import (
	"fmt"
	"strconv"

	"github.com/jake-schuler/fll-event-screen/config"
	"github.com/jake-schuler/fll-event-screen/models"
	"google.golang.org/api/sheets/v4"
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
	var resp *sheets.ValueRange
	var err error

	switch scheduleRound {
	case practice:
		resp, err = config.Srv.Spreadsheets.Values.Get(config.SheetID, "Schedule!A3:E").Do()
	case round1:
		resp, err = config.Srv.Spreadsheets.Values.Get(config.SheetID, "Schedule!F3:J").Do()
	case round2:
		resp, err = config.Srv.Spreadsheets.Values.Get(config.SheetID, "Schedule!K3:O").Do()
	case round3:
		resp, err = config.Srv.Spreadsheets.Values.Get(config.SheetID, "Schedule!P3:T").Do()
	}

	if err != nil {
		fmt.Printf("Error reading schedule for %s: %v\n", ScheduleRound[scheduleRound], err)
		return []models.Match{} // Return empty slice instead of panicking
	}

	var matches []models.Match
	matchMap := make(map[int]*models.Match)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return matches
	}

	for _, row := range resp.Values {
		if len(row) < 5 {
			continue // Skip incomplete rows
		}

		// Parse match ID - handle both string and number formats
		var matchID int
		switch v := row[0].(type) {
		case string:
			matchID, err = strconv.Atoi(v)
			if err != nil {
				fmt.Printf("Error parsing match ID: %v\n", err)
				continue
			}
		case float64:
			matchID = int(v)
		default:
			fmt.Printf("Unexpected match ID type: %T\n", v)
			continue
		}

		// Parse team number - handle both string and number formats
		var teamNum int
		switch v := row[1].(type) {
		case string:
			teamNum, err = strconv.Atoi(v)
			if err != nil {
				fmt.Printf("Error parsing team number: %v\n", err)
				continue
			}
		case float64:
			teamNum = int(v)
		default:
			fmt.Printf("Unexpected team number type: %T\n", v)
			continue
		}

		var teamName string
		if name, ok := row[2].(string); ok {
			teamName = name
		}

		// Parse table assignment
		table, ok := row[4].(string)
		if !ok {
			fmt.Printf("Error parsing table assignment: expected string, got %T\n", row[4])
			continue
		}

		isRedTable := table == "Red A" || table == "Red B"

		// Get or create match
		match, exists := matchMap[matchID]
		if !exists {
			match = &models.Match{
				ID:            matchID,
				IsRedTable:    isRedTable,
				AllianceA:     0,
				AllianceAName: "",
				AllianceB:     0,
				AllianceBName: "",
			}
			matchMap[matchID] = match
		}

		// Assign team to correct alliance based on table
		switch table {
		case "Red A", "Blue A":
			match.AllianceA = teamNum
			match.AllianceAName = teamName
		case "Red B", "Blue B":
			match.AllianceB = teamNum
			match.AllianceBName = teamName
		}
	}

	// Convert map to slice
	for _, match := range matchMap {
		matches = append(matches, *match)
	}

	return matches
}
