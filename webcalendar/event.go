package webcalendar

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

type Event struct {
	Name     string
	StartsAt time.Time
	EndsAt   time.Time

	Address   string
	Building  string
	Classroom string

	Profs      []string
	CFU        int
	TotalHours int

	eventHash string
}

func (e *Event) toIcal() (*ical.Calendar, error) {
	vevent := ical.NewComponent(ical.CompEvent)
	vevent.Props.SetText(ical.PropUID, e.GetHash())
	vevent.Props.SetText(ical.PropSummary, e.GetParsedName())
	vevent.Props.SetText(ical.PropDescription, e.GetParsedDescription())
	vevent.Props.SetText(ical.PropLocation, e.GetParsedLocation())
	vevent.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	vevent.Props.SetDateTime(ical.PropDateTimeStart, e.StartsAt.UTC())
	vevent.Props.SetDateTime(ical.PropDateTimeEnd, e.EndsAt.UTC())

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//Unipi Sync//IT")
	cal.Children = append(cal.Children, vevent)
	return cal, nil
}

func (e *Event) GetHash() string {
	if e.eventHash == "" {
		rawEventId := sha256.Sum256(fmt.Appendf(nil,
			"Random static string|%s|%s|%s|%v|%v",
			e.GetParsedName(),
			e.GetParsedDescription(),
			e.GetParsedLocation(),
			e.StartsAt, e.EndsAt,
		))
		e.eventHash = hex.EncodeToString(rawEventId[:])
	}
	return e.eventHash
}

func (e *Event) GetParsedName() string {
	return fmt.Sprintf("(%s) %s", e.Classroom, e.Name)
}
func (e *Event) GetParsedDescription() string {
	lines := make([]string, 0, 3)

	// Docente: Berry Mr Guy
	if len(e.Profs) > 0 {
		profKey := "Docenti"
		if len(e.Profs) == 1 {
			profKey = "Docente"
		}
		lines = append(lines, fmt.Sprintf("| %s: %s", profKey, strings.Join(e.Profs, ", ")))
	}

	// Durata totale: 90 ore
	if e.TotalHours != 0 {
		lines = append(lines, fmt.Sprintf("| Durata totale: %d ore", e.TotalHours))
	}

	// CFU: 5
	if e.TotalHours != 0 {
		lines = append(lines, fmt.Sprintf("| %d CFU", e.CFU))
	}

	return strings.Join(lines, "\n")
}
func (e *Event) GetParsedLocation() string {
	lines := make([]string, 0, 3)

	// Classroom
	classSplit := strings.Split(e.Classroom, " ")
	if len(classSplit) > 1 {
		lines = append(lines, "Aula "+classSplit[len(classSplit)-1])
	}

	// Actual address
	lines = append(lines, e.Building)
	lines = append(lines, e.Address)

	return strings.Join(lines, ", ")
}
