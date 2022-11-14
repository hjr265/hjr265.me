package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/dslipak/pdf"
)

type Injury struct {
	GameDate      string
	GameTime      string
	Matchup       string
	Team          string
	PlayerName    string
	CurrentStatus string
	Reason        string
}

type State struct {
	LastType   TextType
	ReportDate string
	Injury     Injury
}

func main() {
	log.SetFlags(0)

	// Errors omitted for brevity.
	r, _ := pdf.Open(filename)
	// defer f.Close()

	injuries := []Injury{}
	var state State
	for no := 1; no < r.NumPage(); no++ {
		log.Printf("Page %d", no)

		page := r.Page(no)
		rows, _ := page.GetTextByRow()
		for ri, row := range rows {
			for _, text := range row.Content {
				switch {
				case reInjuryReport.MatchString(text.S):
					state.LastType = TextReportDate
					m := reInjuryReport.FindStringSubmatch(text.S)
					state.ReportDate = m[1]

				case (state.LastType == TextHeader || state.LastType == TextRowBreak) && reHeaders.MatchString(text.S):
					state.LastType = TextHeader

				case (state.LastType == TextHeader || state.LastType == TextRowBreak) && reGameDate.MatchString(text.S):
					state.LastType = TextGameDate
					state.Injury.GameDate = text.S

				case (state.LastType == TextGameDate || state.LastType == TextRowBreak) && reGameTime.MatchString(text.S):
					state.LastType = TextGameTime
					state.Injury.GameTime = text.S

				case state.LastType == TextGameTime && reMatchup.MatchString(text.S):
					state.LastType = TextMatchup
					state.Injury.Matchup = text.S

				case state.LastType == TextMatchup /* && isTeamName(text.S) */ :
					state.LastType = TextTeamName
					state.Injury.Team = text.S

				case state.LastType == TextTeamName /* && isPlayerName(text.S) */ :
					state.LastType = TextPlayerName
					state.Injury.PlayerName = text.S

				case state.LastType == TextRowBreak && text.S != "Game Date" /* && isPlayerName(text.S) */ :
					state.LastType = TextPlayerName
					state.Injury.PlayerName = text.S

				case state.LastType == TextPlayerName /* && isCurrentStatus(text.S) */ :
					state.LastType = TextCurrentStatus
					state.Injury.CurrentStatus = text.S

				case state.LastType == TextCurrentStatus: // Assuming Reason always appears after Current Status.
					state.LastType = TextReason
					state.Injury.Reason = text.S
				}
			}

			if state.Injury.PlayerName != "" && state.LastType == TextReason {
				log.Printf("%2d: %s", ri, state.Injury.PlayerName)
				injuries = append(injuries, state.Injury)
				state.Injury.PlayerName = ""
				state.Injury.CurrentStatus = ""
				state.Injury.Reason = ""
			}

			state.LastType = TextRowBreak
		}

		log.Println()
	}

	b, _ := json.MarshalIndent(injuries, "", "  ")
	os.WriteFile("parsed.json", b, 0755)
}

const filename = "Injury-Report_2022-11-07_07AM.pdf"

var (
	reInjuryReport = regexp.MustCompile(`Injury Report: (\d{2}/\d{2}/\d{2} \d{2}:\d{2} (A|P)M)`)
	reHeaders      = regexp.MustCompile(`^(Game Date|Game Time|Matchup|Team|Player Name|Current Status|Reason)$`)
	reGameDate     = regexp.MustCompile(`\d{2}/\d{2}/\d{2}`)
	reGameTime     = regexp.MustCompile(`\d{2}:\d{2} \([A-Z]{2}\)`)
	reMatchup      = regexp.MustCompile(`[A-Z]{3}@[A-Z]{3}`)
)

type TextType int

const (
	TextUnknown = iota
	TextRowBreak
	TextReportDate
	TextHeader
	TextGameDate
	TextGameTime
	TextMatchup
	TextTeamName
	TextPlayerName
	TextCurrentStatus
	TextReason
)
