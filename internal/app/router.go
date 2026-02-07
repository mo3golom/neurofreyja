package app

import (
	"context"
	"strings"

	"neurofreyja/internal/features/draw_card"
	"neurofreyja/internal/features/find_card_scan"
	"neurofreyja/internal/shared/telegram"

	"gopkg.in/telebot.v3"
)

func (a *App) RegisterHandlers() {
	findService := &find_card_scan.Service{
		Cards:   a.Cards,
		Storage: a.Storage,
		LLM:     a.LLM,
		Model:   a.Config.OpenRouterModelTit,
	}
	findHandler := &find_card_scan.Handler{
		Service:            findService,
		Messenger:          a.Messenger,
		Logger:             a.Logger,
		DeleteAfterMinutes: a.Config.DeleteAfterMinutes,
	}

	drawService := &draw_card.Service{
		Cards:   a.Cards,
		Drawn:   a.Drawn,
		Storage: a.Storage,
		LLM:     a.LLM,
		Model:   a.Config.OpenRouterModelDesc,
	}
	drawHandler := &draw_card.Handler{
		Service:   drawService,
		Messenger: a.Messenger,
		Logger:    a.Logger,
	}

	route := func(c telebot.Context) error {
		msg := c.Message()
		if msg == nil {
			return nil
		}

		text := telegram.MessageText(msg)
		if text == "" {
			return nil
		}

		if telegram.IsGroupChat(msg) {
			if !telegram.BotMentioned(text, a.Config.BotUsername) {
				return nil
			}
		}
		text = telegram.RemoveBotMention(text, a.Config.BotUsername)

		text = strings.TrimSpace(text)
		switch text {
		case "/find_card_scan":
			return findHandler.Handle(c)
		case "/draw_card":
			return drawHandler.Handle(c)
		default:
			_, err := a.Messenger.SendText(context.Background(), msg.Chat, "Такая команда мне неизвестна")
			if err != nil && a.Logger != nil {
				a.Logger.WithError(err).Warn("failed to send unknown command")
			}
			return nil
		}
	}

	a.Bot.Handle(telebot.OnText, route)
	a.Bot.Handle(telebot.OnPhoto, route)
}
