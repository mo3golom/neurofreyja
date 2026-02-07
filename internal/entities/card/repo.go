package card

import "context"

type Repository interface {
	FindByTitles(ctx context.Context, titles []string) ([]Card, error)
	FindRandomNotDrawn(ctx context.Context, chatID int64) (*Card, error)
	UpdateDescription(ctx context.Context, id int64, description string) error
}
