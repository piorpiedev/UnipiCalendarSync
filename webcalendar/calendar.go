package webcalendar

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

func SyncEvents(events []Event, username, passwd, caldavUrl, caldavPath string) error {
	// CalDav client
	cli, err := caldav.NewClient(
		webdav.HTTPClientWithBasicAuth(nil, username, passwd),
		caldavUrl,
	)
	if err != nil {
		return err
	}

	// Fetch remote events
	ctx := context.Background()
	remoteEvents := make(map[string]string) // Map remote UID -> event path
	objects, err := cli.QueryCalendar(ctx, caldavPath, &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{Name: "VCALENDAR"},
	})
	if err != nil {
		return err
	}
	for _, obj := range objects {
		uid := strings.TrimSuffix(obj.Path, ".ics")
		uid = uid[strings.LastIndex(uid, "/")+1:]
		remoteEvents[uid] = obj.Path
	}

	// Parse local events
	localEvents := make(map[string]*Event, len(events))
	for _, e := range events {
		localEvents[e.GetHash()] = &e
	}

	// Upload new events
	for eventHash, localEvent := range localEvents {
		icalEvent, err := localEvent.toIcal()
		if err != nil {
			return err
		}

		eventPath := fmt.Sprintf("%s/%s.ics", caldavPath, eventHash)
		if _, ok := remoteEvents[eventHash]; ok {
			continue
		}

		log.Printf("Creating event: %s", eventHash)
		if _, err := cli.PutCalendarObject(ctx, eventPath, icalEvent); err != nil {
			return err
		}
	}

	// Delete modified or removed events
	for uid, remotePath := range remoteEvents {
		if _, ok := localEvents[uid]; !ok {
			log.Printf("Deleting event: %s", uid)
			if err := cli.RemoveAll(ctx, remotePath); err != nil {
				return err
			}
		}
	}

	return nil

}
