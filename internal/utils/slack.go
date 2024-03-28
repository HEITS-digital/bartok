package utils

import (
	"bartok/internal/constants"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
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
		ImageBlock(constants.WelcomeImageUrl),

		TextToMarkdownBlock(constants.WelcomeIntro1),
		TextToMarkdownBlock(constants.WelcomeIntro2),

		TextToHeaderBlock(constants.HeimsAppHeader),
		TextToMarkdownBlock(constants.HeimsAppSection1),
		TextToMarkdownBlock(constants.HeimsAppSection2),
		TextToMarkdownBlock(constants.HeimsAppSection3),
		TextToMarkdownBlock(constants.HeimsAppSection4),
		TextToMarkdownBlock(constants.HeimsAppSection5),
		TextToMarkdownBlock(constants.HeimsAppSection6),

		TextToHeaderBlock(constants.ChannelsHeader),
		TextToMarkdownBlock(constants.ChannelsSection1),
		TextToMarkdownBlock(constants.ChannelsSection2),

		TextToHeaderBlock(constants.AboutHeader),
		TextToMarkdownBlock(constants.AboutSection1),

		TextToHeaderBlock(constants.OnboardingHeader),
		TextAndOptionsToCheckboxBlock(constants.OnboardingSection1, []string{constants.OnboardingCheckbox1, constants.OnboardingCheckbox2, constants.OnboardingCheckbox3}),
		TextAndOptionsToCheckboxBlock(constants.OnboardingSection2, []string{constants.OnboardingCheckbox4, constants.OnboardingCheckbox5, constants.OnboardingCheckbox6, constants.OnboardingCheckbox7, constants.OnboardingCheckbox8}),
		TextAndOptionsToCheckboxBlock(constants.OnboardingSection3, []string{constants.OnboardingCheckbox9, constants.OnboardingCheckbox10, constants.OnboardingCheckbox11, constants.OnboardingCheckbox12}),

		slack.NewDividerBlock(),

		TextAndOptionsToCheckboxBlock(constants.TrustSection, []string{constants.TrustCheckbox}),
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
func TextAndOptionsToCheckboxBlock(text string, options []string) slack.Block {
	var optionsBlock []*slack.OptionBlockObject
	for index, option := range options {
		// the option value needs to be unique among all the options to be sent
		// If not, options having the same value will be auto-checked when one will be selected
		uuid := uuid.New().String()[:10]
		textBlock := slack.NewTextBlockObject("mrkdwn", option, false, false)
		optionBlock := slack.NewOptionBlockObject(fmt.Sprintf(uuid, index), textBlock, nil)
		optionsBlock = append(optionsBlock, optionBlock)
	}
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil,
		slack.NewAccessory(slack.NewCheckboxGroupsBlockElement(constants.CheckboxActionId, optionsBlock...)),
	)
}
func ImageBlock(imageUrl string) slack.Block {
	return slack.NewImageBlock(imageUrl, "Heits.digital", uuid.New().String()[:5], nil)
}

func quotedText(text string) string {
	return fmt.Sprintf(">*%s*", text)
}
