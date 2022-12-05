package service

import (
	"fmt"
	"github.com/slack-go/slack"
)

type SlackApiService interface {
	SendMessage(channelId string, messageBlocks []slack.Block) error
	SendMessageWithOptions(channelId string, options ...slack.MsgOption) error
	GetAuthUserId() string
}
type slackApiService struct {
	client *slack.Client
}

func (s *slackApiService) SendMessageWithOptions(channelId string, options ...slack.MsgOption) error {
	_, _, err := s.client.PostMessage(channelId, options...)
	return err
}
func (s *slackApiService) SendMessage(channelId string, messageBlocks []slack.Block) error {
	return s.SendMessageWithOptions(channelId, slack.MsgOptionBlocks(messageBlocks...))
}

func NewSlackService(client *slack.Client) SlackApiService {
	return &slackApiService{client: client}
}

func (s *slackApiService) GetAuthUserId() string {
	response, e := s.client.AuthTest()
	if e != nil {
		fmt.Printf("Auth error when trying to get the bot user ID: %s\n", e)
	}
	return response.UserID
}
