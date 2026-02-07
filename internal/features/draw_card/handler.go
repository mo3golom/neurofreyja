package draw_card

import (
	"context"
	"fmt"
	"html"

	"neurofreyja/internal/shared/telegram"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	Service   *Service
	Messenger *telegram.Messenger
	Logger    *logrus.Logger
}

func (h *Handler) Handle(c telebot.Context) error {
	msg := c.Message()
	if msg == nil {
		return nil
	}

	ctx := context.Background()

	loading, err := h.Messenger.SendText(ctx, msg.Chat, "Вытягиваю карту...")
	if err != nil && h.Logger != nil {
		h.Logger.WithError(err).Warn("failed to send loading message")
	}
	defer func() {
		if loading != nil {
			_ = h.Messenger.DeleteMessage(loading)
		}
	}()

	cardItem, description, err := h.Service.Draw(ctx, msg.Chat.ID)
	if err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("draw_card failed")
		}
		return nil
	}
	if cardItem == nil {
		return nil
	}

	image, err := h.Service.FetchImage(ctx, cardItem.ImageID)
	if err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("failed to fetch card image")
		}
		return nil
	}

	escapedTitle := html.EscapeString(cardItem.Title)
	escapedDescription := html.EscapeString(description)
	caption := fmt.Sprintf("Карта \"<b>%s</b>\"\n\n%s", escapedTitle, escapedDescription)
	opts := &telebot.SendOptions{ParseMode: telebot.ModeHTML}
	_, err = h.Messenger.SendPhoto(ctx, msg.Chat, image, caption, opts)
	if err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("failed to send drawn card")
		}
		return nil
	}

	if err := h.Service.MarkDrawn(ctx, cardItem.ID, msg.Chat.ID); err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("failed to mark drawn card")
		}
	}

	return nil
}
