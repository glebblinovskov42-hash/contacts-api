package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type ContactModel struct {
	ID        int        `json:"id,omitempty"`
	Name      string     `json:"name"`
	Number    string     `json:"number"`
	IsFav     bool       `json:"is_favourite"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	return pgx.Connect(ctx, conn)
}

func CreateTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS contacts(
		id SERIAL PRIMARY KEY,
		name VARCHAR(20) NOT NULL,
		number VARCHAR(20) NOT NULL,
		is_favourite BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := conn.Exec(ctx, sqlQuery)

	return err
}

func InsertRaw(ctx context.Context, conn *pgx.Conn, contact ContactModel) error {
	sqlQuery := `
		INSERT INTO contacts (name, number, is_favourite)
		VALUES ($1, $2, $3)
	`

	_, err := conn.Exec(
		ctx,
		sqlQuery,
		contact.Name,
		contact.Number,
		contact.IsFav)

	return err
}

func SelectContacts(ctx context.Context, conn *pgx.Conn) ([]ContactModel, error) {
	sqlQuery := `
		SELECT id, name, number, is_favourite, created_at
		FROM contacts
		ORDER BY id ASC
	`

	contacts := make([]ContactModel, 0)

	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact ContactModel

		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Number,
			&contact.IsFav,
			&contact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)

	}

	return contacts, nil

}

func UpdateContact(ctx context.Context, conn *pgx.Conn, contact ContactModel) error {
	sqlQuery := `
		UPDATE contacts
		SET name=$1, number=$2, is_favourite=$3
		WHERE id=$4;
	`
	_, err := conn.Exec(
		ctx,
		sqlQuery,
		contact.Name,
		contact.Number,
		contact.IsFav,
		contact.ID,
	)
	return err
}

func FavContacts(ctx context.Context, conn *pgx.Conn) ([]ContactModel, error) {
	sqlQuery := `
		SELECT * FROM contacts
		WHERE is_favourite=true
		ORDER BY id ASC
	`

	contacts := make([]ContactModel, 0)

	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact ContactModel

		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Number,
			&contact.IsFav,
			&contact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)

	}

	return contacts, nil
}

func DeleteContact(ctx context.Context, conn *pgx.Conn, id int) error {
	sqlQuery := `
		DELETE FROM contacts
		WHERE id = $1;
	`

	_, err := conn.Exec(ctx, sqlQuery, id)
	return err
}

func ContactExist(ctx context.Context, conn *pgx.Conn, id int) (bool, error) {
	var exist bool
	sqlQuery := `
		SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1);
	`

	err := conn.QueryRow(ctx, sqlQuery, id).Scan(&exist)
	return exist, err
}

func InsertContactWithTx(ctx context.Context, conn *pgx.Conn, contact ContactModel) (int, time.Time, error) {
	var createdAt time.Time
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, createdAt, err
	}

	sqlQuery := `
		INSERT INTO contacts (name, number, is_favourite)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`

	var id int
	rows, err := tx.Query(ctx, sqlQuery, contact.Name, contact.Number, contact.IsFav)
	for rows.Next() {
		rows.Scan(&id, &createdAt)
	}
	defer rows.Close()
	if err != nil {
		tx.Rollback(ctx)
		return 0, createdAt, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, createdAt, err
	}

	return id, createdAt, nil
}
