package history

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type RepoSQLX struct {
	db *sqlx.DB
}

func NewRepoSQLX(db *sqlx.DB) *RepoSQLX {
	return &RepoSQLX{db: db}
}

func (r *RepoSQLX) Insert(ctx context.Context, msg Message) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
insert into neuro_freyja_history
	(message_id, chat_id, chat_title, chat_type, sent_at, content, delete_at)
values ($1, $2, $3, $4, $5, $6, $7)
`

	_, err := r.db.ExecContext(ctx, query,
		msg.MessageID,
		msg.ChatID,
		msg.ChatTitle,
		msg.ChatType,
		msg.SentAt,
		msg.Content,
		msg.DeleteAt,
	)
	return err
}

func (r *RepoSQLX) FindDueDeletions(ctx context.Context) ([]DeletionRecord, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
select id, message_id, chat_id
from neuro_freyja_history
where delete_at is not null
  and delete_at <= now()
  and deleted_at is null
`

	var records []DeletionRecord
	if err := r.db.SelectContext(ctx, &records, query); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *RepoSQLX) MarkDeleted(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
update neuro_freyja_history
set deleted_at = now()
where id = any($1)
`

	_, err := r.db.ExecContext(ctx, query, ids)
	return err
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 5*time.Second)
}
