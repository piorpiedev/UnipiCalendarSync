package cineca

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unipi-calendar-sync/webcalendar"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.Italian, cases.Compact)
var partitionRegex = regexp.MustCompile(`\bcorso \w\b`)

func GetCalendarJson(cinecaCalendarId, cinecaUrl string, from, to time.Time, courseYear int, partition string) ([]webcalendar.Event, error) {
	payload := strings.NewReader(`{
		"mostraImpegniAnnullati":false,
		"mostraIndisponibilitaTotali":false,
		"filtroSfondoId":null,
		"linkCalendarioId":"` + cinecaCalendarId + `",
		"clienteId":"628de8b9b63679f193b87046",
		"pianificazioneTemplate":false,
		"dataInizio":"` + from.Format("2006-01-02T15:04:05") + `.000Z",
		"dataFine":"` + to.Format("2006-01-02T15:04:05") + `.000Z"
	}`) // Cambiato mostraIndisponibilitaTotali:true -> false
	req, err := http.NewRequest("POST", cinecaUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", "__cf_bm="+os.Getenv("CFBM_COOKIE")+"; langUniversityPlanner=it")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// eventHelper := []EventHelper{}
	eventHelper := []map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&eventHelper); err != nil {
		return nil, err
	}
	events := []webcalendar.Event{}

	partitionLower := strings.ToLower(partition)
	thisPartitionRegex := regexp.MustCompile(`\bcorso ` + partitionLower + `\b`)
	for _, eh := range eventHelper {
		detailsMap, ok := getMapField(eh["evento"], "dettagliDidattici")
		if !ok {
			continue
		}
		eventDetails := detailsMap.([]any)[0].(map[string]any)
		eventNameLower := strings.ToLower(eventDetails["nome"].(string))
		eventNameLower = strings.TrimSpace(eventNameLower)

		// Filter just in case
		if courseYear > 0 {
			year := getInt(eventDetails, "annoCorso")
			if year != 0 && year != courseYear {
				continue
			}
		}
		if partition != "" {
			if partitionMap, ok := eventDetails["partizione"]; ok && partitionMap != nil {
				if rawPartition, ok := getMapField(partitionMap, "descrizione"); ok {
					partition := rawPartition.(string)
					partition = strings.TrimPrefix(strings.ToLower(partition), "corso ")
					if partitionLower != partition {
						continue
					}
				}

			} else if partitionRegex.MatchString(eventNameLower) {
				if !thisPartitionRegex.MatchString(eventNameLower) {
					continue
				}
				eventNameLower = strings.Replace(eventNameLower, "corso "+partitionLower, "", 1)
				eventNameLower = strings.TrimSpace(eventNameLower)
			}
		}

		// Parse time
		startsAt, err := time.Parse("2006-01-02T15:04:05Z", eh["dataInizio"].(string))
		if err != nil {
			continue
		}
		endsAt, err := time.Parse("2006-01-02T15:04:05Z", eh["dataFine"].(string))
		if err != nil {
			continue
		}

		event := webcalendar.Event{
			Name:       strings.ToUpper(eventNameLower),
			CFU:        getInt(eventDetails, "cfu"),
			TotalHours: getInt(eventDetails, "totaleOre"),
			StartsAt:   startsAt,
			EndsAt:     endsAt,
		}

		// Add prof, description, and so on
		resources := eh["risorse"].([]any)
		for _, res := range resources {
			if profMap, ok := getMapField(res, "docente"); ok {
				prof := profMap.(map[string]any)
				profName := prof["nome"].(string) + " " + prof["cognome"].(string)
				event.Profs = append(event.Profs, titleCaser.String(profName))
				continue
			}
			if classroomMap, ok := getMapField(res, "aula"); ok {
				classroom := classroomMap.(map[string]any)
				event.Classroom = classroom["codice"].(string)

				if buildingMap, ok := classroom["edificio"]; ok && buildingMap != nil {
					building := buildingMap.(map[string]any)
					event.Building = building["descrizione"].(string)
					event.Address = fmt.Sprintf("%s n.%s (%s)", titleCaser.String(building["via"].(string)), building["numeroCivico"], building["comune"])
				}
				continue
			}
		}

		events = append(events, event)
	}

	return events, nil
}

func getInt(from map[string]any, key string) int {
	if v, ok := from[key]; ok && v != nil {
		return int(v.(float64))
	}
	return 0
}
func getMapField(from any, key string) (any, bool) {
	v, ok := from.(map[string]any)[key]
	return v, ok && v != nil
}
