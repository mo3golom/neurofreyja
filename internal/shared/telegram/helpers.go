package telegram

import (
	"strings"

	"gopkg.in/telebot.v3"
)

func MessageText(msg *telebot.Message) string {
	if msg == nil {
		return ""
	}
	if msg.Caption != "" {
		return msg.Caption
	}
	return msg.Text
}

func ChatType(msg *telebot.Message) string {
	if msg == nil || msg.Chat == nil {
		return ""
	}
	return string(msg.Chat.Type)
}

func IsGroupChat(msg *telebot.Message) bool {
	if msg == nil || msg.Chat == nil {
		return false
	}
	return msg.Chat.Type == telebot.ChatGroup || msg.Chat.Type == telebot.ChatSuperGroup
}

func NormalizeBotUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.TrimPrefix(username, "@")
	return username
}
