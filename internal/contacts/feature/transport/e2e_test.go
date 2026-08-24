package transport

import (
	"bytes"
	service "contacts-api/internal/contacts/feature/service"
	postgres "contacts-api/internal/contacts/feature/storage"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const testDBConn = "postgres://postgres:sos11982@localhost:5433/postgres_test?sslmode=disable"

func TestCreateContacte2e(t *testing.T) {
	logger := zap.NewNop()

	conn, err := pgx.Connect(context.Background(), testDBConn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS contacts (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			number TEXT NOT NULL,
			is_favourite BOOLEAN DEFAULT false, 
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = conn.Exec(context.Background(), "TRUNCATE contacts RESTART IDENTITY")
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	storage := postgres.NewStorage(conn)
	service := service.NewContactService(storage, logger)
	handler := NewHandler(service, logger)

	router := mux.NewRouter()
	router.HandleFunc("/contacts", handler.CreateContact).Methods("POST")

	body := []byte(`{"name":"E2ETest", "num":"12345678911", "isfav":true}`)
	req := httptest.NewRequest("POST", "/contacts", bytes.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	t.Logf("Response body: %s", w.Body.String())

	if w.Code != http.StatusCreated {
		t.Errorf("expecred status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	if _, ok := resp["ID"]; !ok {
		t.Errorf("expecred ID in response")
	}
}
