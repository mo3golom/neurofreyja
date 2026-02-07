package drawn

import "context"

type Repository interface {
	Insert(ctx context.Context, cardID int64, chatID int64) error
}
