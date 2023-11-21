package service

import (
	"bartok/internal/constants"
	"bartok/internal/datastruct"
	"bartok/internal/utils"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type GoogleApiService interface {
	GetGoogleCalendarService() GoogleCalendarService
}

type googleApiService struct {
	client *http.Client
}

type googleCalendarService struct {
	service *calendar.Service
}

type GoogleCalendarService interface {
	GetEmployeeBirthdays(timeNow time.Time) []datastruct.EmployeeEvent
}

func NewGoogleApiService() GoogleApiService {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
	}
	ctx := context.Background()

	token := new(oauth2.Token)

	token.RefreshToken = os.Getenv("GOOGLE_REFRESH_TOKEN")

	return &googleApiService{client: config.Client(ctx, token)}
}

func NewGoogleCalendarService(client *http.Client) GoogleCalendarService {
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}

	return &googleCalendarService{service: service}
}

func (google *googleApiService) GetGoogleCalendarService() GoogleCalendarService {
	return NewGoogleCalendarService(google.client)
}

func (googleCalendar *googleCalendarService) AttachCalendar(calendarId string) *calendar.CalendarListEntry {
	entry := new(calendar.CalendarListEntry)
	entry.Id = calendarId

	googleCalendar.service.CalendarList.Insert(entry)

	return entry
}

func GenerateBirthdayWish(description string) string {
	expression, _ := regexp.Compile("\\<@(.*?)\\>")

	if !expression.MatchString(description) {
		return description
	}

	matches := strings.Join(expression.FindAllString(description, -1), "")

	// when the description contains only the slack user id we want to generate a random message
	if matches == description {
		return fmt.Sprintf(utils.GetRandomItem(constants.RandomBirthdayWishes), matches)
	}

	return description
}

func CreateEmployeeBirthdayEvent(item *calendar.Event, date string) datastruct.EmployeeEvent {
	var employeeEvent datastruct.EmployeeEvent
	employeeEvent.Id = item.Id
	employeeEvent.Date = date
	employeeEvent.Employee = item.Summary
	employeeEvent.Text = GenerateBirthdayWish(html.UnescapeString(item.Description))
	employeeEvent.Type = datastruct.Birthday

	return employeeEvent
}

func (googleCalendar *googleCalendarService) GetEmployeeBirthdays(timeNow time.Time) []datastruct.EmployeeEvent {
	calendarId := os.Getenv("GOOGLE_CALENDAR_EMPLOYEES")
	googleCalendar.AttachCalendar(calendarId)
	const (
		YYYYMMDD = "2006-01-02"
	)
	utcTime := timeNow.UTC()

	startTime := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), 0, 0, 0, 0, utcTime.Location()).Format(time.RFC3339)
	endTime := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), 23, 59, 59, 0, utcTime.Location()).Format(time.RFC3339)

	var formattedEvents = make([]datastruct.EmployeeEvent, 0)

	events, _ := googleCalendar.service.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(startTime).TimeMax(endTime).Do()

	if len(events.Items) == 0 || events == nil {
		return nil
	} else {
		for _, item := range events.Items {
			isBirthdayEvent, _ := regexp.MatchString(string(datastruct.Birthday), item.Summary)
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}

			if isBirthdayEvent && utcTime.Format(YYYYMMDD) == date {
				formattedEvents = append(formattedEvents, CreateEmployeeBirthdayEvent(item, date))
			}
		}
	}

	return formattedEvents
}
