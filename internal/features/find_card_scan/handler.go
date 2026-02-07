package find_card_scan

import (
	"context"

	"neurofreyja/internal/shared/telegram"
	ftime "neurofreyja/internal/shared/time"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	Service            *Service
	Messenger          *telegram.Messenger
	Logger             *logrus.Logger
	DeleteAfterMinutes int
}

func (h *Handler) Handle(c telebot.Context) error {
	msg := c.Message()
	if msg == nil {
		return nil
	}

	ctx := context.Background()

	loading, err := h.Messenger.SendText(ctx, msg.Chat, "Сейчас поищу...")
	if err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("failed to send loading message")
		}
	}
	defer func() {
		if loading != nil {
			_ = h.Messenger.DeleteMessage(loading)
		}
	}()

	if msg.ReplyTo == nil {
		_, err := h.Messenger.SendText(ctx, msg.Chat, "Эту команду надо вызывать реплаем на сообщение")
		if err != nil && h.Logger != nil {
			h.Logger.WithError(err).Warn("failed to send reply requirement message")
		}
		return nil
	}

	replyText := telegram.MessageText(msg.ReplyTo)
	cards, err := h.Service.FindCards(ctx, replyText)
	if err != nil {
		if h.Logger != nil {
			h.Logger.WithError(err).Warn("find_card_scan failed")
		}
		return nil
	}
	if len(cards) == 0 {
		return nil
	}

	for _, card := range cards {
		image, err := h.Service.FetchImage(ctx, card.ImageID)
		if err != nil {
			if h.Logger != nil {
				h.Logger.WithError(err).Warn("failed to fetch card image")
			}
			continue
		}

		deleteAt := ftime.DeleteAtMinutes(h.DeleteAfterMinutes)

		_, err = h.Messenger.SendPhotoWithDelete(ctx, msg.Chat, image, "Это карта - \""+card.Title+"\"", deleteAt)
		if err != nil {
			if h.Logger != nil {
				h.Logger.WithError(err).Warn("failed to send card photo")
			}
			continue
		}

		_, err = h.Messenger.SendTextWithDelete(ctx, msg.Chat, "У вас есть 10 минут чтобы скачать изображения", deleteAt)
		if err != nil {
			if h.Logger != nil {
				h.Logger.WithError(err).Warn("failed to send delete warning")
			}
		}
	}

	return nil
}
