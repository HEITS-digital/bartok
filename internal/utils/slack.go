package utils

import (
	"bartok/internal/constants"
	"fmt"
	"github.com/slack-go/slack"
	"regexp"
)

func GetRandomWelcomeMessage(user string) string {
	return fmt.Sprintf(GetRandomItem(constants.WelcomeMessages), user)
}

func GetFarewellMessage(user string, joke string) string {
	return fmt.Sprintf(constants.FarewellMessageFormat, user, joke)
}

func GetNewMemberDM() []slack.Block {

	blocks := []slack.Block{
		TextToHeaderBlock(constants.WelcomeHeader),
		TextToMarkdownBlock(constants.WelcomeIntro),
		TextToMarkdownBlock(constants.WelcomeInternalMeetings),
		slack.NewDividerBlock(),
		TextToMarkdownBlock(constants.WelcomeSlackChannels),
		slack.NewDividerBlock(),
		TextToMarkdownBlock(fmt.Sprintf(constants.WelcomePto,
			//Florina Condulet
			"U01GSDPJSGH",
			//Andreea Caraba
			"U028QJ7R339",
		)),
		slack.NewDividerBlock(),
		TextToMarkdownBlock(constants.WelcomeStayConnected),
		slack.NewDividerBlock(),
		TextToMarkdownBlock(fmt.Sprintf(constants.WelcomeOutro,
			//Teodora Cenan
			"U01TLCC4TRC",
		)),
	}
	return blocks
}

func GetRandomReply(user string, messages []string) string {
	return fmt.Sprintf(GetRandomItem(messages), user)
}

func RemoveMentionFromText(text string) string {
	// in order to process a received test, we'll get rid of the mention part inside
	reg := regexp.MustCompile(`<([^)]+)>`)
	return reg.ReplaceAllString(text, "")
}
func TextToHeaderBlock(text string) slack.Block {
	return slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", text, false, false))

}
func TextToTextBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("plain_text", text, true, false),
		nil,
		nil,
	)
}
func TextToQuotedTextBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", quotedText(text), false, false),
		nil,
		nil,
	)
}
func TextToMarkdownBlock(text string) slack.Block {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil,
		nil,
	)
}

func quotedText(text string) string {
	return fmt.Sprintf(">*%s*", text)
}
