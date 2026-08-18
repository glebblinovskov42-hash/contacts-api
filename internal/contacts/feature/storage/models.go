package postgres

import "time"

type ContactModel struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Number    string     `db:"number"`
	IsFav     bool       `db:"is_favourite"`
	CreatedAt *time.Time `db:"created_at,omitempty"`
}
