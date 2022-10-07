package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/abadojack/whatlanggo"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
var botUserId string

func main() {
	http.HandleFunc("/ask", slashCommandHandler)

	http.HandleFunc("/slack/actions", func(w http.ResponseWriter, r *http.Request) {
		var payload slack.InteractionCallback
		err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
		if err != nil {
			fmt.Printf("Could not parse action response JSON: %v", err)
		}
		for _, action := range payload.ActionCallback.BlockActions {
			switch action.ActionID {
			case "another_question_action":
				handleAnotherQuestion(payload)
			default:
				fmt.Printf(`seems like unhandled action with ID: %s`, action.ActionID)
			}
		}

	})

	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("am primit event")
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
			eventHandler(eventsAPIEvent)
		}
	})

	http.HandleFunc("/cron/watercooler", watercoolerHandler)
	port := os.Getenv("PORT")
	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleAnotherQuestion(payload slack.InteractionCallback) {
	if payload.Message.ReplyCount > 0 {
		fmt.Println("Am intrat aici, e bine!")
	}
}

func watercoolerAction() {

	channelId := os.Getenv("WATERCOOLER_CHANNEL_ID")
	blocks := createMessageBlocksForWaterCooler()
	_, _, err := api.PostMessage(channelId, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		log.Fatal(err)
	}
}
func watercoolerHandler(_ http.ResponseWriter, _ *http.Request) {
	watercoolerAction()
}

func createMessageBlocksForWaterCooler() []slack.Block {
	waterCoolerEnIntro := fmt.Sprintf(
		"%s\n%s",
		getRandomMessage(waterCoolerGreetings),
		getRandomMessage(waterCoolerEnIntros),
	)
	waterCoolerRoIntro := getRandomMessage(waterCoolerRoIntros)

	question := getRandomFromMap(waterCoolerQuestions)
	waterCoolerEnQuestion := fmt.Sprintf(">*%s*", question[0])
	waterCoolerRoQuestion := fmt.Sprintf(">*%s*", question[1])

	var blocks []slack.Block
	blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", waterCoolerEnIntro, true, false), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", waterCoolerEnQuestion, false, false), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", waterCoolerRoIntro, true, false), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", waterCoolerRoQuestion, false, false), nil, nil))
	return blocks
}

type SlashCommandResponse struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Text     string `json:"text"`
	ImageUrl string `json:"image_url"`
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	var data SlashCommandResponse
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/chucknorris":
		data = SlashCommandResponse{
			ResponseType: "in_channel",
			Text:         ":chucknorris: " + getChuckNorrisJoke(),
		}
	case "/truth":
		params := &slack.Msg{Text: s.Text}
		answer, err := getAnswer()
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.HasSuffix(params.Text, "?") {
			data = SlashCommandResponse{
				ResponseType: "in_channel",
				Attachments:  []Attachment{{Text: strings.ToUpper(answer.Answer), ImageUrl: answer.Image}},
			}
		} else {
			answer.Answer = "This doesn't seem to be a question. Try harder!"
			data = SlashCommandResponse{
				ResponseType: "ephemeral",
				Text:         answer.Answer,
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func eventHandler(eventsAPIEvent slackevents.EventsAPIEvent) {
	innerEvent := eventsAPIEvent.InnerEvent
	fmt.Println("Incoming callback event: ", innerEvent)

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		// detect the language of the received text
		info := whatlanggo.Detect(removeMentionFromText(ev.Text))
		var message string
		//verify if this is a question
		if strings.HasSuffix(ev.Text, "?") {
			if info.Lang == whatlanggo.Eng {
				message = getRandomReply(ev.User, randomEnAnswers)
			} else {
				// if no english is detected, just reply in romanian
				message = getRandomReply(ev.User, randomRoAnswers)
			}
		} else {
			if info.Lang == whatlanggo.Eng {
				message = getRandomReply(ev.User, randomEnReplies)
			} else {
				message = getRandomReply(ev.User, randomRoReplies)
			}
		}
		// verify if this mention comes from a thread and reply back if so
		if len(ev.ThreadTimeStamp) > 0 {
			api.PostMessage(ev.Channel, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true), slack.MsgOptionTS(ev.ThreadTimeStamp))
		} else {
			postMessage(*api, ev.Channel, message)
		}

	case *slackevents.MemberJoinedChannelEvent:
		if !(len(botUserId) > 0) {
			botUserId = getBotUserId(*api)
		}
		// send messages when a non bot user is added and when added to a specified channel id
		if ev.User != botUserId && ev.Channel == os.Getenv("GENERAL_CHANNEL_ID") {
			postMessage(*api, ev.Channel, getRadomWelcomeMessage(ev.User))
			postMessage(*api, ev.User, getNewMemberDM())
		}
	}
}

func verifyRequestAndRespond(w http.ResponseWriter, body []byte) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(r.Challenge))
}

func removeMentionFromText(text string) string {
	// in order to process a received test, we'll going to get rid of the mention part inside
	reg := regexp.MustCompile(`\<([^)]+)\>`)
	return reg.ReplaceAllString(text, "")
}

func getRandomFromMap(mapItems [][]string) []string {
	rand.Seed(time.Now().Unix())
	index := rand.Int() % len(mapItems)
	return mapItems[index]
}

func getRandomMessage(messages []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(messages)
	return messages[n]
}

func getRandomReply(user string, messages []string) string {
	return fmt.Sprintf(getRandomMessage(messages), user)
}

func getBotUserId(api slack.Client) string {
	response, e := api.AuthTest()
	if e != nil {
		fmt.Printf("Auth error when trying to get the bot user ID: %s\n", e)
	}
	return response.UserID
}

func postMessage(api slack.Client, channel, message string) string {
	respChannel, _, err := api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("Error sending slack message: %s\n", err)
	}
	return respChannel
}
