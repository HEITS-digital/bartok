package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		fmt.Println("Missing SLACK_TOKEN in environment")
		os.Exit(1)
	}
	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	// messages := 0
	var botChannelJoinedEventReceived bool

	for msg := range rtm.IncomingEvents {
		// fmt.Print("Event Received: ")

		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.MemberJoinedChannelEvent:
			fmt.Printf("Member joined information: %v\n", ev)
			channel := ev.Channel
			user := ev.User
			if !botChannelJoinedEventReceived {
				postMessage(*api, channel, getRadomWelcomeMessage(user))
				postMessage(*api, user, getNewMemberDM())
			}

		case *slack.MemberLeftChannelEvent:
			fmt.Printf("Member left information: %v\n", ev)
			channel := ev.Channel
			postMessage(*api, channel, "Farewell amigo! :wave:\nWe're really going to miss trying to avoid you around here")

		case *slack.ConnectedEvent:
			// Ignore connected

		case *slack.ChannelJoinedEvent:
			// this flag would be overwritten from the send message
			botChannelJoinedEventReceived = true

		case *slack.DesktopNotificationEvent:
			if ev.IsChannelInvite && botChannelJoinedEventReceived {
				channel := ev.Channel
				postMessage(*api, channel, "Hi all! I'm Bartók the goat.\nI'm still trying to figure out stuff so be pacient and don't ping me with messages for now.")
				botChannelJoinedEventReceived = false
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			//fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
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
	return fmt.Sprintf(teamJoinWelcomeMessageFormat, "CSKGXKXS5", "G01GE16SBAP", "CEC2Y6QD9", "C01S8NR19TR", "C01NY7FN34Y")
}

const teamJoinWelcomeMessageFormat string = `Welcome to HEITS.digital :wave: ! We are super excited that you joined us, and wish you the best of luck on this new adventure. 
I’m Bartók the goat, and I am here to share some useful information with you:
*1. Internal meetings*
- Each Monday at 11am we have the Internal & Informal meeting, where we discuss important company updates.
- Once a month we meet and share knowledge, during the HEITS talks initiative. Come and find out cool stuff, both technical and non-technical.
*2. Slack channels*
- If you ever need help from our workspace’s administrators, please reach out in #general
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
