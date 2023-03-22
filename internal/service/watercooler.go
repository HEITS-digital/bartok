package service

import (
	"bartok/internal/constants"
	"bartok/internal/datastruct"
	"bartok/internal/repository"
	"bartok/internal/utils"
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
		utils.TextToTextBlock(greeting),
		utils.TextToTextBlock(enIntro),
		utils.TextToQuotedTextBlock(question.English),
		utils.TextToTextBlock(roIntro),
		utils.TextToQuotedTextBlock(question.Romanian),
	}
	return blocks
}
