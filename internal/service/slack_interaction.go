package service

import (
	"bartok/internal/client"
	"bartok/internal/constants"
	"bartok/internal/datastruct"
	"bartok/internal/utils"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"os"
	"strings"
)

type SlackInteractionService interface {
	SlashCommands(slack.SlashCommand) (*datastruct.SlashCommandResponse, error)
	SlackEvents(eventsAPIEvent slackevents.EventsAPIEvent) error
}
type slackInteractionService struct {
	slackService SlackApiService
	botUserId    *string
}

func NewSlackInteractionService(slackService SlackApiService) SlackInteractionService {
	return &slackInteractionService{slackService: slackService}
}

func (s *slackInteractionService) SlashCommands(command slack.SlashCommand) (*datastruct.SlashCommandResponse, error) {
	var data datastruct.SlashCommandResponse

	switch command.Command {
	case "/chucknorris":
		data = datastruct.SlashCommandResponse{
			ResponseType: "in_channel",
			Text:         ":chucknorris: " + client.GetChuckNorrisJoke(),
		}
	case "/truth":
		params := &slack.Msg{Text: command.Text}
		answer, err := client.GetAnswer()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if strings.HasSuffix(params.Text, "?") {
			data = datastruct.SlashCommandResponse{
				ResponseType: "in_channel",
				Attachments:  []datastruct.Attachment{{Text: strings.ToUpper(answer.Answer), ImageUrl: answer.Image}},
			}
		} else {
			answer.Answer = "This doesn't seem to be a question. Try harder!"
			data = datastruct.SlashCommandResponse{
				ResponseType: "ephemeral",
				Text:         answer.Answer,
			}
		}
	}
	return &data, nil
}

func (s *slackInteractionService) SlackEvents(eventsAPIEvent slackevents.EventsAPIEvent) error {
	innerEvent := eventsAPIEvent.InnerEvent
	fmt.Println("Incoming callback event: ", innerEvent)

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		// detect the language of the received text
		info := whatlanggo.Detect(utils.RemoveMentionFromText(ev.Text))
		var message string
		//verify if this is a question
		if strings.HasSuffix(ev.Text, "?") {
			if info.Lang == whatlanggo.Eng {
				message = utils.GetRandomReply(ev.User, constants.RandomEnAnswers)
			} else {
				// if no english is detected, just reply in romanian
				message = utils.GetRandomReply(ev.User, constants.RandomRoAnswers)
			}
		} else {
			if info.Lang == whatlanggo.Eng {
				message = utils.GetRandomReply(ev.User, constants.RandomEnReplies)
			} else {
				message = utils.GetRandomReply(ev.User, constants.RandomRoReplies)
			}
		}
		// verify if this mention comes from a thread and reply if so
		if len(ev.ThreadTimeStamp) > 0 {
			_ = s.slackService.SendMessageWithOptions(ev.Channel, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true), slack.MsgOptionTS(ev.ThreadTimeStamp))
		} else {
			_ = s.postMessage(ev.Channel, message)
		}

	case *slackevents.MemberJoinedChannelEvent:
		// send messages when a non bot user is added and when added to a specified channel id
		if ev.User != *s.botUserId && ev.Channel == os.Getenv("GENERAL_CHANNEL_ID") {
			go s.postMessage(ev.Channel, utils.GetRandomWelcomeMessage(ev.User))
			go s.postMessage(ev.User, utils.GetNewMemberDM())
		}
	}
	return nil
}

func (s *slackInteractionService) getBotUserId() string {
	if s.botUserId == nil {
		var botUserId string
		botUserId = s.slackService.GetAuthUserId()
		s.botUserId = &botUserId
	}
	return *s.botUserId

}

func (s *slackInteractionService) postMessage(channel, message string) error {
	err := s.slackService.SendMessageWithOptions(channel, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("Error sending slack message: %s\n", err)
	}
	return nil
}
