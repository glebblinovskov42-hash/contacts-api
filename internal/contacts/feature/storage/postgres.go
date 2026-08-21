package postgres

import (
	"contacts-api/internal/core/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(conn *pgx.Conn) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) GetAll(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	sqlQuery := `
		SELECT id, name, number, is_favourite, created_at
		FROM contacts
		ORDER by id ASC
		LIMIT $1 OFFSET $2;
	`

	rows, err := s.conn.Query(ctx, sqlQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contacts []domain.Contact
	for rows.Next() {
		var c domain.Contact
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Number,
			&c.IsFav,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	return contacts, nil
}

func (s *Storage) Save(ctx context.Context, contact domain.Contact) (int, time.Time, error) {
	var id int
	var createdAt time.Time

	sqlQuery := `
		INSERT INTO contacts (name, number, is_favourite)
		VALUES ($1, $2, $3)
		RETURNING id, Created_at
	`

	err := s.conn.QueryRow(ctx, sqlQuery, contact.Name, contact.Number, contact.IsFav).Scan(&id, &createdAt)
	if err != nil {
		return 0, time.Time{}, err
	}
	return id, createdAt, nil
}

func (s *Storage) GetById(ctx context.Context, id int) (domain.Contact, error) {
	var c domain.Contact
	sqlQuery := `
		SELECT id, name, number, is_favourite, created_at
		FROM contacts
		WHERE id = $1;
	`

	err := s.conn.QueryRow(ctx, sqlQuery, id).Scan(
		&c.ID,
		&c.Name,
		&c.Number,
		&c.IsFav,
		&c.CreatedAt,
	)
	if err != nil {
		return domain.Contact{}, err
	}
	return c, nil
}

func (s *Storage) Update(ctx context.Context, contact domain.Contact) error {
	sqlQuery := `
		UPDATE contacts
		SET name=$1, number=$2, is_favourite=$3
		WHERE id=$4
	`

	_, err := s.conn.Exec(ctx, sqlQuery, contact.Name, contact.Number, contact.IsFav, contact.ID)
	return err
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	sqlQuery := `
		DELETE FROM contacts WHERE id=$1;
	`

	_, err := s.conn.Exec(ctx, sqlQuery, id)

	return err
}

func (s *Storage) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	sqlQuery := `
		SELECT EXISTS(SELECT 1 FROM contacts WHERE id=$1);
	`
	err := s.conn.QueryRow(ctx, sqlQuery, id).Scan(&exists)
	return exists, err
}
