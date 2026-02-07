package card

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type RepoSQLX struct {
	db *sqlx.DB
}

func NewRepoSQLX(db *sqlx.DB) *RepoSQLX {
	return &RepoSQLX{db: db}
}

func (r *RepoSQLX) FindByTitles(ctx context.Context, titles []string) ([]Card, error) {
	if len(titles) == 0 {
		return nil, nil
	}

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
select title, image_id
from neuro_freyja_card
where lower(title) = any($1)
`

	var cards []Card
	if err := r.db.SelectContext(ctx, &cards, query, titles); err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *RepoSQLX) FindRandomNotDrawn(ctx context.Context, chatID int64) (*Card, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
select c.id, c.title, c.image_id, c.description
from neuro_freyja_card c
where c.id not in (
	select card_id from neuro_freyja_drawn_card where chat_id=$1
)
order by random()
limit 1
`

	var card Card
	if err := r.db.GetContext(ctx, &card, query, chatID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

func (r *RepoSQLX) UpdateDescription(ctx context.Context, id int64, description string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	const query = `
update neuro_freyja_card
set description=$1, updated_at=now()
where id=$2
`

	_, err := r.db.ExecContext(ctx, query, description, id)
	return err
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 5*time.Second)
}
