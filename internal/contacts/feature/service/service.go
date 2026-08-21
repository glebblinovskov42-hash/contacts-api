package service

import (
	"contacts-api/internal/core/domain"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrContactsNotFound = errors.New("contacts not found")
	ErrContactNodFound  = errors.New("contact not found")
	ErrInvalidContact   = errors.New("invalid contact data")
)

type ContactStorage interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Contact, error)
	Save(ctx context.Context, contact domain.Contact) (int, time.Time, error)
	GetById(ctx context.Context, id int) (domain.Contact, error)
	Update(ctx context.Context, contact domain.Contact) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, id int) (bool, error)
}

type ContactService struct {
	storage ContactStorage
	logger  *zap.Logger
}

func NewContactService(storage ContactStorage, logger *zap.Logger) *ContactService {
	return &ContactService{
		storage: storage,
		logger:  logger,
	}
}

func lenValidation(name, number string) error {
	if name == "" {
		return ErrInvalidContact
	}
	if len(number) != 11 {
		return ErrInvalidContact
	}
	return nil
}

func (s *ContactService) GetAllContacts(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	contacts, err := s.storage.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, ErrContactsNotFound
	}

	return contacts, nil
}

func (s *ContactService) GetContactById(ctx context.Context, id int) (domain.Contact, error) {
	contact, err := s.storage.GetById(ctx, id)
	if err != nil {
		return domain.Contact{}, ErrContactNodFound
	}
	return contact, nil
}

func (s *ContactService) CreateContact(ctx context.Context, name, number string, isFav bool) (domain.Contact, error) {
	lenValidation(name, number)

	contact := domain.Contact{
		Name:   name,
		Number: number,
		IsFav:  isFav,
	}

	id, createdAt, err := s.storage.Save(ctx, contact)
	if err != nil {
		return domain.Contact{}, err
	}

	contact.ID = id
	contact.CreatedAt = &createdAt
	return contact, nil
}

func (s *ContactService) UpdateContact(ctx context.Context, contact domain.Contact) error {
	lenValidation(contact.Name, contact.Number)

	exists, err := s.storage.Exists(ctx, contact.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrContactNodFound
	}

	return s.storage.Update(ctx, contact)
}

func (s *ContactService) DeleteContact(ctx context.Context, id int) error {
	exists, err := s.storage.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrContactNodFound
	}

	return s.storage.Delete(ctx, id)
}

func (s *ContactService) GetFavContacts(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	all, err := s.storage.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var favs []domain.Contact
	for _, c := range all {
		if c.IsFav {
			favs = append(favs, c)
		}
	}

	if len(favs) == 0 {
		return nil, ErrContactNodFound
	}

	return favs, nil
}
