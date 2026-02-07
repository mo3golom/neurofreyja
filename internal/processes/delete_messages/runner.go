package delete_messages

import (
	"context"
	"strings"
	"time"

	"neurofreyja/internal/entities/history"
	"neurofreyja/internal/shared/telegram"

	"github.com/sirupsen/logrus"
)

type Runner struct {
	History   history.Repository
	Messenger *telegram.Messenger
	Logger    *logrus.Logger
	Interval  time.Duration
}

func (r *Runner) Run(ctx context.Context) {
	if r.Interval <= 0 {
		r.Interval = time.Minute
	}

	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.sweep(ctx)
		}
	}
}

func (r *Runner) sweep(ctx context.Context) {
	records, err := r.History.FindDueDeletions(ctx)
	if err != nil {
		if r.Logger != nil {
			r.Logger.WithError(err).Warn("failed to fetch deletions")
		}
		return
	}
	if len(records) == 0 {
		return
	}

	var deleted []int64
	for _, record := range records {
		err := r.Messenger.DeleteByID(record.ChatID, record.MessageID)
		if err != nil {
			if isTerminalDeleteError(err) {
				if r.Logger != nil {
					r.Logger.WithError(err).WithFields(logrus.Fields{
						"history_id": record.ID,
						"chat_id":    record.ChatID,
						"message_id": record.MessageID,
					}).Info("message already deleted or cannot be deleted")
				}
				deleted = append(deleted, record.ID)
				continue
			}
			if r.Logger != nil {
				r.Logger.WithError(err).WithFields(logrus.Fields{
					"history_id": record.ID,
					"chat_id":    record.ChatID,
					"message_id": record.MessageID,
				}).Warn("failed to delete message")
			}
			continue
		}
		deleted = append(deleted, record.ID)
	}

	if err := r.History.MarkDeleted(ctx, deleted); err != nil {
		if r.Logger != nil {
			r.Logger.WithError(err).Warn("failed to mark deletions")
		}
	}
}

func isTerminalDeleteError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "message to delete not found"),
		strings.Contains(msg, "message can't be deleted"),
		strings.Contains(msg, "message cannot be deleted"),
		strings.Contains(msg, "message_id_invalid"),
		strings.Contains(msg, "message id invalid"):
		return true
	default:
		return false
	}
}
