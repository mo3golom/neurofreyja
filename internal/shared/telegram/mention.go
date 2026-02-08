package telegram

import (
	"regexp"
	"strings"
)

func BotMentioned(text, username string) bool {
	username = NormalizeBotUsername(username)
	if username == "" {
		return false
	}
	re := mentionRegex(username)
	return re.MatchString(text)
}

func RemoveBotMention(text, username string) string {
	username = NormalizeBotUsername(username)
	if username == "" {
		return strings.TrimSpace(text)
	}
	re := mentionRegex(username)
	cleaned := re.ReplaceAllString(text, "$1")
	return strings.TrimSpace(cleaned)
}

func mentionRegex(username string) *regexp.Regexp {
	pattern := `(?i)(^|[^\w]|/\w+)@` + regexp.QuoteMeta(username) + `(\b|$)`
	return regexp.MustCompile(pattern)
}
