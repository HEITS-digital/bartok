package internal

import (
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

type HttpServer interface {
	Start(port string) error
}
type httpServer struct {
	slackInteractionService service.SlackInteractionService
	watercoolerService      service.WatercoolerService
	birthdayService         service.HappyBirthdayService
}

func CreateHttpServer(
	slackInteractionService service.SlackInteractionService,
	watercoolerService service.WatercoolerService,
	birthdayService service.HappyBirthdayService,
) HttpServer {
	return &httpServer{
		slackInteractionService: slackInteractionService,
		watercoolerService:      watercoolerService,
		birthdayService:         birthdayService,
	}
}

func (h *httpServer) Start(port string) error {
	fmt.Printf("Server listening on port %s\n", port)

	http.HandleFunc("/status", h.statusHandler)
	http.HandleFunc("/ask", h.slackAskHandler)
	http.HandleFunc("/slack/events", h.slackEventsHandler)
	http.HandleFunc("/cron/watercooler", h.watercoolerHandler)
	http.HandleFunc("/cron/birthday", h.birthdayHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (h *httpServer) slackEventsHandler(w http.ResponseWriter, r *http.Request) {

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
		err := h.slackInteractionService.SlackEvents(eventsAPIEvent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *httpServer) slackAskHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.slackInteractionService.SlashCommands(s)

	response, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (h *httpServer) watercoolerHandler(http.ResponseWriter, *http.Request) {
	err := h.watercoolerService.PostNewQuestion()
	if err != nil {
		return
	}
}

func (h *httpServer) birthdayHandler(w http.ResponseWriter, r *http.Request) {
	cards := h.birthdayService.PostBirthDayCards()

	response, err := json.Marshal(cards)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)

}

func (h *httpServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := make(map[string]string)
	response["message"] = "ok"
	json, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return
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
