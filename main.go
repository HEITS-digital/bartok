package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			verifyRequestAndRespond(w, body)
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			handleSlackCallbackEvents(eventsAPIEvent)
		}
	})
	fmt.Println("Server listening")
	http.ListenAndServe(":8080", nil)
}

func handleSlackCallbackEvents(eventsAPIEvent slackevents.EventsAPIEvent) {
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
		// avoid sending messages when the bot is added to a channel
		if ev.User != botUserId {
			postMessage(*api, ev.Channel, getRadomWelcomeMessage(ev.User))
			postMessage(*api, ev.User, getNewMemberDM())
		}

	case *slackevents.MemberLeftChannelEvent:
		postMessage(*api, ev.Channel, "Farewell amigo! :wave:\nWe're really going to miss trying to avoid you around here.")
	}
}

func removeMentionFromText(text string) string {
	// in order to process a received test, we'll going to get rid of the mention part inside
	reg := regexp.MustCompile(`\<([^)]+)\>`)
	return reg.ReplaceAllString(text, "")
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

func getRandomReply(user string, messages []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(messages)
	return fmt.Sprintf(messages[n], user)
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

func getRadomWelcomeMessage(user string) string {
	rand.Seed(time.Now().Unix())
	messages := []string{
		"Hello <@%v>, and welcome to HEITS :wave:",
		"HOORAY!\nWelcome to the team <@%v> :dog_hooray:",
		"Ciao <@%v>\nBenvenuto nella nostra famiglia! :mafia:",
	}
	n := rand.Int() % len(messages)
	return fmt.Sprintf(messages[n], user)
}

func getNewMemberDM() string {
	return fmt.Sprintf(teamJoinWelcomeMessageFormat, "CEC0Z16QL", "CSKGXKXS5", "C02054LCV6E", "CEC2Y6QD9", "C01S8NR19TR", "C01NY7FN34Y")
}

var randomEnReplies = []string{
	"Hey <@%v>. Where is the beef?",
	"Sorry <@%v> but I can't deal with you now.\nThis week is so very busy and my skin is broken",
	"Yes <@%v>\nI have superpowers because I was born at a very young age",
	"Stand back <@%v>, your hair makes me nervous",
	"Hey <@%v>.\nWould you like to kiss my flamingo? :flamingo:",
	"<@%v> on a scale of 1 to 5, how anxious are you when using public bathrooms?",
	"Stop asking for my number <@%v>!!!",
	"Are you afraid of raccoons <@%v>?",
	"Pickled cabbage -> that's my secret\nWhat's yours <@%v>?",
	"<@%v> -> :talktothehand:",
	"<@%v> -> :lalala:",
}

var randomRoReplies = []string{
	"Hai sa lasam prostiile pt alta data <@%v>",
	"<@%v>, pt binele tuturor hai sa pretindem ca ti-ai vazut de treaba ta",
	"Oare care-l platiti pe <@%v>? Sa va rog sa-i dati ceva de lucru",
	"Hai sa continuam mai tarziu <@%v>. Acum am de finalizat o comanda la ikea",
	"<@%v> ai observat ca nu-ti raspund in private s-acum incerci aici?",
	"N-am timp acum <@%v>. Hai sa ne auzim mai tarziu... mult mai tarziu",
	"<@%v> -> :talktothehand:",
	"<@%v> -> :lalala:",
}

var randomEnAnswers = []string{
	"Is that a serious question <@%v>?",
	"Hey <@%v> you know what?\nI’ll answer you in a bit. I’m now waiting for motivation to build up",
	"Sorry <@%v> but I just found something important to deal with at this moment",
	"I have no idea how to answer this <@%v> :thinking_face:",
	"I don't think I'm qualified to answer this now <@%v>",
	"Only questions and questions... You know, I have questions too <@%v>. But who's curious of it?",
}

var randomRoAnswers = []string{
	"Mda… Alta intrebare <@%v>! :alta-intrebare:",
	"Habar n-am ce sa-ti raspund la asta <@%v>. Lasa-ma sa ma mai gandesc",
	"Haha <@%v>. Ce te face sa crezi ca am timp pt intrebari acum?",
	"Revin imediat cu un raspuns <@%v>. Momentan mi-am luat o pauza pt gustare :leafy_green",
	"Scuze <@%v>, dar momentan sunt lucruri mai importante de facut decat sa caut raspunsuri pentru tine",
	"Ia zi <@%v>, cine crezi ca te plateste sa ma iei pe mine la intrebari?",
	"<@%v> ai observat ca nu-ti raspund in private s-acum incerci aici? :grin:",
}

const teamJoinWelcomeMessageFormat string = `Welcome to HEITS.digital :wave: ! We are super excited that you joined us, and wish you the best of luck on this new adventure. 
I’m Bartók the goat, and I am here to share some useful information with you:
*1. Internal meetings*
- Each Monday at 11am we have the Internal & Informal meeting, where we disc	uss important company updates.
- Once a month we meet and share knowledge, during the HEITS talks initiative. Come and find out cool stuff, both technical and non-technical.
*2. Slack channels*
- If you ever need help from our workspace’s administrators, please reach out in <#%s>
- Engineering -> <#%s>
- Administrative & Financial stuff -> <#%s>
- Games, Hobbies & Fun -> <#%s>, <#%s> & <#%s>
There are quite a few other channels, depending on your interests or location. Just click on the :heavy_plus_sign: next to the channel list in the sidebar, and click Browse Channels to search for anything that interests you.
*3. PTO*
- This is our vacation planner https://heims.heits.digital/. You can authenticate using your Google account and add your vacation days here. Your Google calendar will later reflect the PTO days.
- For any other information regarding our benefits, or other administrative aspects, you can always reach Lidia Rusu from HR or Florina Condulet from Finance & Administration.
*4. Stay connected*
- Here’s our website https://heits.digital/ - check it out
- Facebook page: https://www.facebook.com/heits.digital - Like & Share
- Linkedin page: https://www.linkedin.com/company/heits-digital/ - Follow & Share
Hope I could be of help and I am working on adding new useful functions. If you have any suggestions, please drop a message to the engineering team.
Sit back, relax, enjoy our community and have fun! :happygoat:`
