package main

import (
	"bartok/internal"
	"bartok/internal/repository"
	"bartok/internal/service"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {

	slackClient := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
	firestoreClient, err := repository.NewFirestoreClient(os.Getenv("FIRESTORE_PROJECT"))
	if err != nil {
		log.Panicln(err)
		return
	}
	dao := repository.NewDAO(firestoreClient)
	slackService := service.NewSlackService(slackClient)

	server := internal.CreateHttpServer(
		service.NewSlackInteractionService(slackService),
		service.NewWatercoolerService(slackService, dao),
		service.NewHappyBirthdayService(slackService),
	)

	_ = server.Start(os.Getenv("PORT"))
}
