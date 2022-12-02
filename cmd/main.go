package main

import (
	"bartok/internal/repository"
	"bartok/internal/service"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"log"
	"net/http"
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
	slackInteractionService := service.NewSlackInteractionService(slackService)
	watercoolerService := service.NewWatercoolerService(slackService, dao)

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Am intrat la /ask")
		s, err := slack.SlashCommandParse(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, err := slackInteractionService.SlashCommands(s)

		response, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(response)
	})
	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			verifyRequestAndRespond(w, body)
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			err := slackInteractionService.SlackEvents(eventsAPIEvent)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})

	http.HandleFunc("/cron/watercooler", func(writer http.ResponseWriter, request *http.Request) {
		err := watercoolerService.PostNewQuestion()
		if err != nil {
			return
		}
	})
	port := os.Getenv("PORT")
	fmt.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func verifyRequestAndRespond(w http.ResponseWriter, body []byte) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal(body, &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(r.Challenge))
}
