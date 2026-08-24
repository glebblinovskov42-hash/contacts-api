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
	ErrContactNotFound  = errors.New("contact not found")
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

func lenValidation(logger *zap.Logger, name, number string) error {
	if name == "" {
		logger.Warn("Empty name")
		return ErrInvalidContact
	}
	if len(number) != 11 {
		logger.Warn("Wrong number")
		return ErrInvalidContact
	}
	return nil
}

func (s *ContactService) GetAllContacts(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	s.logger.Info("Get all contacts", zap.Int("limit", limit), zap.Int("offset", offset))
	contacts, err := s.storage.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return nil, err
	}
	if len(contacts) == 0 {
		s.logger.Warn("Have 0 contacts")
		return nil, ErrContactsNotFound
	}

	return contacts, nil
}

func (s *ContactService) GetContactById(ctx context.Context, id int) (domain.Contact, error) {
	s.logger.Info("Get contact by id", zap.Int("id", id))
	contact, err := s.storage.GetById(ctx, id)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return domain.Contact{}, ErrContactNotFound
	}
	return contact, nil
}

func (s *ContactService) CreateContact(ctx context.Context, name, number string, isFav bool) (domain.Contact, error) {
	s.logger.Info("Create contact", zap.String("Name", name), zap.String("Number", number), zap.Bool("is favourite", isFav))
	lenValidation(s.logger, name, number)

	contact := domain.Contact{
		Name:   name,
		Number: number,
		IsFav:  isFav,
	}

	id, createdAt, err := s.storage.Save(ctx, contact)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return domain.Contact{}, err
	}

	contact.ID = id
	contact.CreatedAt = &createdAt
	return contact, nil
}

func (s *ContactService) UpdateContact(ctx context.Context, contact domain.Contact) error {
	s.logger.Info("Update contact", zap.Int("id", contact.ID))
	lenValidation(s.logger, contact.Name, contact.Number)

	exists, err := s.storage.Exists(ctx, contact.ID)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return err
	}
	if !exists {
		s.logger.Warn("Contact not exist", zap.Error(ErrContactNotFound))
		return ErrContactNotFound
	}

	return s.storage.Update(ctx, contact)
}

func (s *ContactService) DeleteContact(ctx context.Context, id int) error {
	s.logger.Info("Delete contact", zap.Int("id", id))
	exists, err := s.storage.Exists(ctx, id)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return err
	}
	if !exists {
		s.logger.Warn("Contact not exist", zap.Error(ErrContactNotFound))
		return ErrContactNotFound
	}

	return s.storage.Delete(ctx, id)
}

func (s *ContactService) GetFavContacts(ctx context.Context, limit, offset int) ([]domain.Contact, error) {
	s.logger.Info("get favourite contacts", zap.Int("limit", limit), zap.Int("offset", offset))
	all, err := s.storage.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("internal error", zap.Error(err))
		return nil, err
	}

	var favs []domain.Contact
	for _, c := range all {
		if c.IsFav {
			favs = append(favs, c)
		}
	}

	if len(favs) == 0 {
		s.logger.Warn("have 0 favourite contacts", zap.Error(ErrContactNotFound))
		return nil, ErrContactsNotFound
	}

	return favs, nil
}
