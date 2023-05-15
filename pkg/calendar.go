package pkg

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Calendar struct {
	Config  *oauth2.Config
	Log     *logrus.Logger
	service *calendar.Service
}

func (c *Calendar) GenerateUrl(start_date, end_date string, schoolid int) string {
	return c.Config.AuthCodeURL("state-token", oauth2.AccessTypeOnline) + url.QueryEscape(fmt.Sprintf("?%s+07:00?%s+07:00?%d", start_date, end_date, schoolid))
}

func (c *Calendar) NewService(auth string) *Calendar {

	tok, err := c.Config.Exchange(context.TODO(), auth)
	if err != nil {
		c.Log.Errorf("Unable to retrieve token from web: %v", err)
	}
	client := c.Config.Client(context.Background(), tok)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.Log.Errorf("Unable to retrieve Calendar client: %v", err)
	}
	c.service = srv
	return c
}
func (c *Calendar) Create(start_date, end_date, schname string) string {
	event := &calendar.Event{
		Summary: "Pengenalan Sekolah " + schname,
		Start: &calendar.EventDateTime{
			DateTime: start_date,
			TimeZone: "Asia/Singapore",
		},
		End: &calendar.EventDateTime{
			DateTime: end_date,
			TimeZone: "Asia/Singapore",
		},

		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: "Pengenalan Sekolah " + schname,
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=1"},
	}

	calendarId := "primary"
	event, err := c.service.Events.Insert(calendarId, event).ConferenceDataVersion(1).Do()
	if err != nil {
		c.Log.Errorf("Unable to create event: %v", err)
	}
	return event.ConferenceData.EntryPoints[0].Uri
}
