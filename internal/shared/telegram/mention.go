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
	pattern := `(?i)@` + regexp.QuoteMeta(username)
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

func RemoveBotMention(text, username string) string {
	username = NormalizeBotUsername(username)
	if username == "" {
		return strings.TrimSpace(text)
	}
	pattern := `(?i)@` + regexp.QuoteMeta(username)
	re := regexp.MustCompile(pattern)
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}
