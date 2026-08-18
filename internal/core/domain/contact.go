package domain

import "time"

type Contact struct {
	ID        int
	Name      string
	Number    string
	IsFav     bool
	CreatedAt *time.Time
}
