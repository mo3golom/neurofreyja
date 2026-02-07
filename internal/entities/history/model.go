package history

import "time"

type Message struct {
	ID        int64      `db:"id"`
	MessageID int        `db:"message_id"`
	ChatID    int64      `db:"chat_id"`
	ChatTitle string     `db:"chat_title"`
	ChatType  string     `db:"chat_type"`
	SentAt    time.Time  `db:"sent_at"`
	Content   string     `db:"content"`
	DeleteAt  *time.Time `db:"delete_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type DeletionRecord struct {
	ID        int64 `db:"id"`
	MessageID int   `db:"message_id"`
	ChatID    int64 `db:"chat_id"`
}
