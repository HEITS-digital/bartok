package service

import (
	"bartok/internal/datastruct"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/slack-go/slack"
)

type HappyBirthdayService interface {
	PostBirthDayCards() []datastruct.EmployeeEvent
}
type happyBirthdayService struct {
	slackService SlackApiService
}

func (w *happyBirthdayService) PostBirthDayCards() []datastruct.EmployeeEvent {
	events := NewGoogleApiService().GetGoogleCalendarService().GetEmployeeBirthdays(time.Now())
	channelId := os.Getenv("BIRTHDAY_CHANNEL_ID")
	var sentCards = make([]datastruct.EmployeeEvent, 0)
	converter := md.NewConverter("", true, nil)

	for _, event := range events {
		if len(event.Text) > 0 {
			markdown, _ := converter.ConvertString(event.Text)
			if len(markdown) > 0 {
				markdown = strings.ReplaceAll(markdown, `\_`, `_`)
			}
			err := w.slackService.SendMessage(channelId, []slack.Block{
				textToTextWithMrkdwnBlock(markdown),
			})
			event.IsSent = err == nil
		} else {
			event.IsSent = false
		}

		sentCards = append(sentCards, event)
	}

	return sentCards
}

func NewHappyBirthdayService(slackService SlackApiService) HappyBirthdayService {
	return &happyBirthdayService{slackService: slackService}
}

func textToTextWithMrkdwnBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil,
		nil,
	)
}
