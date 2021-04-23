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
		fmt.Print("Event Received: ")

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
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

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
	return fmt.Sprintf(teamJoinWelcomeMessageFormat, "CEC0Z16QL", "CPBLBS3SL", "CSKGXKXS5", "G01GE16SBAP", "C01S8NR19TR", "C01NY7FN34Y")
}

const teamJoinWelcomeMessageFormat = `Welcome to the HEITS Slack Workspace! I'm Bartók the goat, and hopefully I'll have new functions available soon. :crossed_fingers:
If you ever need help from our workspace's administrators, please reach out in <#%s>.
Here's our website <https://heits.digital/>, which you might want to check it out as well. Any website related stuff is discussed here <#%s>.
This is our vacation planner <https://heims.heits.digital/>. You can authenticate using your Google account and add your vacation days inside so you calendar would reflect the PTO days. 
Here's a list of a few other channels you could join:
- Engineering -> <#%s>
- Administrative & Financial stuff -> <#%s>
- Games & Hobbies -> <#%s> & <#%s>

There are quite a few other channels, depending on your interests or location. Just click on the :heavy_plus_sign: next to the channel list in the sidebar, and click Browse Channels to search for anything that interests you.
Now, enjoy our community and have fun! :happygoat:`
