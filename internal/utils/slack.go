package utils

import (
	"bartok/internal/constants"
	"fmt"
	"regexp"
)

func GetRandomWelcomeMessage(user string) string {
	return fmt.Sprintf(GetRandomItem(constants.WelcomeMessages), user)
}

func GetFarewellMessage(user string, joke string) string {
	return fmt.Sprintf(constants.FarewellMessageFormat, user, joke)
}

func GetNewMemberDM() string {
	return fmt.Sprintf(constants.TeamJoinWelcomeMessageFormat, "CEC0Z16QL", "CSKGXKXS5", "C02054LCV6E", "CEC2Y6QD9", "C01S8NR19TR", "C01NY7FN34Y")
}

func GetRandomReply(user string, messages []string) string {
	return fmt.Sprintf(GetRandomItem(messages), user)
}

func RemoveMentionFromText(text string) string {
	// in order to process a received test, we'll get rid of the mention part inside
	reg := regexp.MustCompile(`<([^)]+)>`)
	return reg.ReplaceAllString(text, "")
}
