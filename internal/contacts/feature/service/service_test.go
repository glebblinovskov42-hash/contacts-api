package service

import (
	"contacts-api/internal/core/domain"
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLenValidation(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name     string
		contact  domain.Contact
		hasError bool
	}{
		{"valid contact", domain.Contact{Name: "Bob", Number: "12345678911"}, false},
		{"empty name", domain.Contact{Name: "", Number: "12345678911"}, true},
		{"wrong number len", domain.Contact{Name: "Bob", Number: "123456789112"}, true},
		{"wrong number len and empty name", domain.Contact{Name: "", Number: "123456789112"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lenValidation(logger, tt.contact.Name, tt.contact.Number)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}

		})
	}
}

type MockStorage struct {
	saveFunc func(contact domain.Contact) (int, time.Time, error)
}

func (m *MockStorage) GetAll(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	return nil, nil
}

func (m *MockStorage) Save(ctx context.Context, contact domain.Contact) (int, time.Time, error) {
	return m.saveFunc(contact)
}

func (m *MockStorage) GetById(ctx context.Context, id int) (domain.Contact, error) {
	return domain.Contact{}, nil
}

func (m *MockStorage) Update(ctx context.Context, contact domain.Contact) error {
	return nil
}

func (m *MockStorage) Delete(ctx context.Context, id int) error {
	return nil
}

func (m *MockStorage) Exists(ctx context.Context, id int) (bool, error) {
	return false, nil
}

func TestCreateContact(t *testing.T) {
	logger := zap.NewNop()
	storage := &MockStorage{
		saveFunc: func(contact domain.Contact) (int, time.Time, error) {
			return 1, time.Now(), nil
		},
	}

	service := NewContactService(storage, logger)

	contact, err := service.CreateContact(context.Background(), "Alice", "12345678911", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if contact.Name != "Alice" {
		t.Errorf("expected Alice, got: %s", contact.Name)
	}
}
