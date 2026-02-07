package telegram

import (
	"bytes"
	"context"
	"time"

	"neurofreyja/internal/entities/history"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type Messenger struct {
	Bot     *telebot.Bot
	Logger  *logrus.Logger
	History history.Repository
}

func NewMessenger(bot *telebot.Bot, logger *logrus.Logger, historyRepo history.Repository) *Messenger {
	return &Messenger{
		Bot:     bot,
		Logger:  logger,
		History: historyRepo,
	}
}

func (m *Messenger) SendText(ctx context.Context, chat *telebot.Chat, text string, opts ...interface{}) (*telebot.Message, error) {
	msg, err := m.Bot.Send(chat, text, opts...)
	if err != nil {
		return nil, err
	}
	m.logMessage(ctx, msg, nil)
	return msg, nil
}

func (m *Messenger) SendTextWithDelete(ctx context.Context, chat *telebot.Chat, text string, deleteAt time.Time, opts ...interface{}) (*telebot.Message, error) {
	msg, err := m.Bot.Send(chat, text, opts...)
	if err != nil {
		return nil, err
	}
	m.logMessage(ctx, msg, &deleteAt)
	return msg, nil
}

func (m *Messenger) SendPhoto(ctx context.Context, chat *telebot.Chat, data []byte, caption string, opts ...interface{}) (*telebot.Message, error) {
	msg, err := m.sendPhoto(chat, data, caption, opts...)
	if err != nil {
		return nil, err
	}
	m.logMessage(ctx, msg, nil)
	return msg, nil
}

func (m *Messenger) SendPhotoWithDelete(ctx context.Context, chat *telebot.Chat, data []byte, caption string, deleteAt time.Time, opts ...interface{}) (*telebot.Message, error) {
	msg, err := m.sendPhoto(chat, data, caption, opts...)
	if err != nil {
		return nil, err
	}
	m.logMessage(ctx, msg, &deleteAt)
	return msg, nil
}

func (m *Messenger) DeleteMessage(msg *telebot.Message) error {
	if msg == nil {
		return nil
	}
	return m.Bot.Delete(msg)
}

func (m *Messenger) DeleteByID(chatID int64, messageID int) error {
	return m.Bot.Delete(&telebot.Message{
		ID:   messageID,
		Chat: &telebot.Chat{ID: chatID},
	})
}

func (m *Messenger) sendPhoto(chat *telebot.Chat, data []byte, caption string, opts ...interface{}) (*telebot.Message, error) {
	photo := &telebot.Photo{
		File: telebot.File{
			FileReader: bytes.NewReader(data),
		},
		Caption: caption,
	}

	return m.Bot.Send(chat, photo, opts...)
}

func (m *Messenger) logMessage(ctx context.Context, msg *telebot.Message, deleteAt *time.Time) {
	if m.History == nil || msg == nil {
		return
	}

	record := history.Message{
		MessageID: msg.ID,
		ChatID:    msg.Chat.ID,
		ChatTitle: msg.Chat.Title,
		ChatType:  string(msg.Chat.Type),
		SentAt:    msg.Time().UTC(),
		Content:   MessageText(msg),
		DeleteAt:  deleteAt,
	}

	if err := m.History.Insert(ctx, record); err != nil {
		if m.Logger != nil {
			m.Logger.WithError(err).Warn("failed to log outgoing message")
		}
	}
}
