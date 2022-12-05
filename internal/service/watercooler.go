package service

import (
	"bartok/internal/constants"
	"bartok/internal/datastruct"
	"bartok/internal/repository"
	"bartok/internal/utils"
	"fmt"
	"github.com/slack-go/slack"
	"os"
)

type WatercoolerService interface {
	PostNewQuestion() error
}
type watercoolerService struct {
	slackService SlackApiService
	dao          repository.DAO
}

func (w *watercoolerService) PostNewQuestion() error {
	channelId := os.Getenv("WATERCOOLER_CHANNEL_ID")
	question, err := w.dao.NewWatercoolerQuery().GetNextUnreadQuestion()
	if err != nil {
		return err
	}
	messageBlocks := createMessageBlocksForWatercooler(question)
	err = w.slackService.SendMessage(channelId, messageBlocks)
	if err != nil {
		return err
	}
	question.IsRead = true
	return w.dao.NewWatercoolerQuery().UpdateQuestion(question)

}

func NewWatercoolerService(slackService SlackApiService, dao repository.DAO) WatercoolerService {
	return &watercoolerService{dao: dao, slackService: slackService}
}

func createMessageBlocksForWatercooler(question *datastruct.Question) []slack.Block {
	greeting := utils.GetRandomItem(constants.WatercoolerGreetings)
	enIntro := utils.GetRandomItem(constants.WatercoolerEnIntros)
	roIntro := utils.GetRandomItem(constants.WatercoolerRoIntros)

	blocks := []slack.Block{
		textToTextBlock(greeting),
		textToTextBlock(enIntro),
		textToQuotedTextBlock(question.English),
		textToTextBlock(roIntro),
		textToQuotedTextBlock(question.Romanian),
	}
	return blocks
}

func textToTextBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("plain_text", text, true, false),
		nil,
		nil,
	)
}
func textToQuotedTextBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", quotedText(text), false, false),
		nil,
		nil,
	)
}

func quotedText(text string) string {
	return fmt.Sprintf(">*%s*", text)
}
