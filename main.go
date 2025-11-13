package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"unipi-calendar-sync/cineca"
	"unipi-calendar-sync/webcalendar"
)

const Day = time.Hour * 24

type CalendarConfig struct {
	CinecaUrl        string `json:"cinecaUrl"`
	CinecaCalendarId string `json:"cinecaCalendarId"`
	CalDavUrl        string `json:"caldavUrl"`
	CalDavPath       string `json:"caldavPath"`
	Year             int    `json:"year"`      // Ignored if <= 0
	Partition        string `json:"partition"` // Letter | Ignored if empty
}

func main() {
	// Load calendars config
	jsonConfig, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	calendars := []CalendarConfig{}
	json.NewDecoder(jsonConfig).Decode(&calendars)

	// Sync every calendar
	for _, calendar := range calendars {
		events, err := cineca.GetCalendarJson(
			calendar.CinecaCalendarId,
			calendar.CinecaUrl,
			time.Now().UTC().Add(-Day*365), // Basically maxing out, since the calendar owner already takes care of cleaning up old events. Thank- you-
			time.Now().UTC().Add(Day*365),
			calendar.Year,
			calendar.Partition,
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Discovered %d events\n", len(events))

		// Upload and delete
		if err := webcalendar.SyncEvents(
			events, os.Getenv("CALDAV_USERNAME"), os.Getenv("CALDAV_PASSWD"),
			calendar.CalDavUrl, calendar.CalDavPath); err != nil {
			log.Fatal(err)
		}
	}

}
