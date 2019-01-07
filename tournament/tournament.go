//Package tournament contains the solution for the exercism Tournament exercise
package tournament

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

type teamRecord struct {
	Name   string
	Played int
	Won    int
	Drawn  int
	Lost   int
	Points int
}

func addResult(teams map[string]*teamRecord, line string) error {
	if strings.HasPrefix(line, "#") {
		return nil
	}
	if strings.HasPrefix(line, "\n") {
		_, i := utf8.DecodeRuneInString("\n")
		line = line[i:]
	}
	if line == "" {
		return nil
	}
	row := strings.Split(line, ";")
	if len(row) != 3 {
		return errors.New("invalid data")
	}
	teamA := row[0]
	teamB := row[1]
	a, ok := teams[teamA]
	if !ok {
		a = &teamRecord{Name: teamA}
		teams[teamA] = a
	}
	b, ok := teams[teamB]
	if !ok {
		b = &teamRecord{Name: teamB}
		teams[teamB] = b
	}
	a.Played++
	b.Played++
	switch row[2] {
	case "win":
		a.Won++
		a.Points += 3
		b.Lost++
	case "loss":
		b.Won++
		b.Points += 3
		a.Lost++
	case "draw":
		a.Drawn++
		a.Points++
		b.Drawn++
		b.Points++
	default:
		return fmt.Errorf("bad result: %q", row[2])
	}
	return nil
}

//Tally will convert a semicolon seperated list to a table of tournament results
func Tally(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	teams := make(map[string]*teamRecord)
	for scanner.Scan() {
		if err := addResult(teams, scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	var teamSlice []*teamRecord
	for _, value := range teams {
		teamSlice = append(teamSlice, value)
	}
	sort.Slice((teamSlice)[:], func(i, j int) bool {
		if teamSlice[i].Points == teamSlice[j].Points {
			return teamSlice[i].Name < teamSlice[j].Name
		}
		return teamSlice[i].Points > teamSlice[j].Points
	})
	if _, err := fmt.Fprintf(writer,"%-31v| MP |  W |  D |  L |  P\n", "Team"); err != nil {
		return err
	}
	for _, record := range teamSlice {
		if _, err := fmt.Fprintf(writer,"%-31v| %2v | %2v | %2v | %2v | %2v\n", record.Name, record.Played, record.Won, record.Drawn, record.Lost, record.Points); err != nil {
			return err
		}
	}
	return nil
}