package telegram

import (
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
