package app

import (
	"context"

	"neurofreyja/internal/features/draw_card"
	"neurofreyja/internal/features/find_card_scan"

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

	a.Bot.Handle("/find_card_scan", func(c telebot.Context) error {
		return findHandler.Handle(c)
	})
	a.Bot.Handle("/draw_card", func(c telebot.Context) error {
		return drawHandler.Handle(c)
	})
	a.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		_, err := a.Messenger.SendText(context.Background(), c.Chat(), "Такая команда мне неизвестна")
		if err != nil && a.Logger != nil {
			a.Logger.WithError(err).Warn("failed to send unknown command")
		}
		return nil
	})
}
