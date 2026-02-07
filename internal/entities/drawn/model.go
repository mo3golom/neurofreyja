package drawn

type Drawn struct {
	ID     int64 `db:"id"`
	CardID int64 `db:"card_id"`
	ChatID int64 `db:"chat_id"`
}
