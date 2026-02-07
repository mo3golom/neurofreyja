package drawn

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

func (r *RepoSQLX) Insert(ctx context.Context, cardID int64, chatID int64) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `insert into neuro_freyja_drawn_card(card_id, chat_id) values ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, cardID, chatID)
	return err
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 5*time.Second)
}
