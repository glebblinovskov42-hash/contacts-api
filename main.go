package main

import (
	"api/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

/*
	type contact struct {
		Name   string
		Number string
		IsFav  bool
	}

var allContacts []contact
*/
var contact db.ContactModel
var mtx = sync.RWMutex{}
var ctx = context.Background()
var conn, err = db.CreateConnection(ctx)

func numValidation(contact db.ContactModel) error {
	if contact.Name == "" {
		return errors.New("Передано пустое поле")
	} else if len(contact.Number) != 11 {
		return errors.New("Передан не валидный номер")
	}
	return nil
}

func handleCreateContact(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()

	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Ошибка чтения", http.StatusBadRequest)
		return
	}

	if err := numValidation(contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.InsertRaw(ctx, conn, contact); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalIndent(contact, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func handleAllContacts(w http.ResponseWriter, r *http.Request) {
	defer mtx.RUnlock()
	mtx.RLock()

	contacts, err := db.SelectContacts(ctx, conn)
	if err != nil {
		panic(err)
	}

	if len(contacts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()

	idStr := mux.Vars(r)["ID"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}

	exist, err := db.ContactExist(ctx, conn, id)
	if err != nil {
		panic(err)
	}

	if !exist {
		http.Error(w, "Контакт не найден", http.StatusNotFound)
		return
	}

	contact.ID = id

	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Ошибка чтения", http.StatusBadRequest)
		return
	}

	if err := numValidation(contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := db.UpdateContact(ctx, conn, contact); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contact, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func handleFavContacts(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()

	contacts, err := db.FavContacts(ctx, conn)
	if err != nil {
		panic(err)
	}

	if len(contacts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()

	ID := mux.Vars(r)["ID"]
	id, err := strconv.Atoi(ID)
	if err != nil {
		panic(err)
	}

	exist, err := db.ContactExist(ctx, conn, id)
	if err != nil {
		panic(err)
	}

	if !exist {
		http.Error(w, "Контакт не найден", http.StatusNotFound)
		return
	}

	if err := db.DeleteContact(ctx, conn, id); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	if err := db.CreateTable(ctx, conn); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Path("/contacts").Methods("POST").HandlerFunc(handleCreateContact)
	router.Path("/contacts").Methods("GET").HandlerFunc(handleAllContacts)
	router.Path("/contacts/{ID}").Methods("PUT").HandlerFunc(handleUpdateContact)
	router.Path("/contacts/fav").Methods("GET").HandlerFunc(handleFavContacts)
	router.Path("/contacts/{ID}").Methods("DELETE").HandlerFunc(handleDeleteContact)

	if err := http.ListenAndServe(":9091", router); err != nil {
		fmt.Println("Ошибка сервера")
		return
	}
}
