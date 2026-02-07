package history

import "context"

type Repository interface {
	Insert(ctx context.Context, msg Message) error
	FindDueDeletions(ctx context.Context) ([]DeletionRecord, error)
	MarkDeleted(ctx context.Context, ids []int64) error
}
