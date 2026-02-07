package card

import "database/sql"

type Card struct {
	ID          int64          `db:"id"`
	Title       string         `db:"title"`
	ImageID     string         `db:"image_id"`
	Description sql.NullString `db:"description"`
}
