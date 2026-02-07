package delete_messages

import (
	"context"
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
