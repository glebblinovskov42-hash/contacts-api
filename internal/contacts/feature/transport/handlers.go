package transport

import (
	service "contacts-api/internal/contacts/feature/service"
	"contacts-api/internal/core/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ContactRequest struct {
	Name   string `json:"name"`
	Number string `json:"num"`
	IsFav  bool   `json:"isfav"`
}

type ContactResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Number    string     `json:"number"`
	IsFav     bool       `json:"is_favourite"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Handler struct {
	service *service.ContactService
	logger  *zap.Logger
}

func NewHandler(service *service.ContactService, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("CreateContact called", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Request body decoded", zap.Any("req", req))

	contact, err := h.service.CreateContact(r.Context(), req.Name, req.Number, req.IsFav)
	if err != nil {
		h.logger.Warn("Failed to create contact", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("Contact created", zap.Int("id", contact.ID), zap.String("name", contact.Name))

	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalIndent(contact, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func (h *Handler) AllContacts(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("AllContacts called", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	contacts, err := h.service.GetAllContacts(r.Context())
	if err != nil {
		if errors.Is(err, service.ErrContactsNotFound) {
			h.logger.Warn("Have 0 contacts", zap.Error(err))
			http.Error(w, "contacts not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("AllContacts completed", zap.Int("count", len(contacts)), zap.Int("status", http.StatusOK))

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func (h *Handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("UpdateContact called", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		h.logger.Warn("Invalid id", zap.Error(err))
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	contact := domain.Contact{
		ID:     id,
		Name:   req.Name,
		Number: req.Number,
		IsFav:  req.IsFav,
	}

	if err := h.service.UpdateContact(r.Context(), contact); err != nil {
		if errors.Is(err, service.ErrContactNodFound) {
			h.logger.Warn("Contact not found", zap.Error(err))
			http.Error(w, "contact not found", http.StatusNotFound)
			return
		}
		h.logger.Warn("Bad request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("UpdateContact completed", zap.Int("id", contact.ID), zap.Int("status", http.StatusOK))

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contact, "", "    ")
	if err != nil {
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("DeleteContact called", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		h.logger.Warn("Invalid id", zap.Error(err))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteContact(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrContactNodFound) {
			h.logger.Warn("Contact not found", zap.Error(err))
			http.Error(w, "contact not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("DeleteContact completed", zap.Int("id", id), zap.Int("status", http.StatusNoContent))

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) FavContacts(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("FavContacts called", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	contacts, err := h.service.GetFavContacts(r.Context())
	if err != nil {
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(contacts) == 0 {
		h.logger.Warn("Have 0 favourites contacts", zap.Error(err))
		http.Error(w, "no fav contacts", http.StatusNotFound)
		return
	}

	h.logger.Info("FavContacts completed", zap.Int("count", len(contacts)), zap.Int("status", http.StatusOK))

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		h.logger.Error("Internal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}
