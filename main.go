package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

/*
	func Map[T, U any](s []T, fn func(T) U) []U {
		var res []U
		for _, v := range s{
			res = append(res, fn(v))
		}
		return res
	}
*/
type contact struct {
	Name   string
	Number string
	IsFav  bool
}

var allContacts []contact
var mtx = sync.RWMutex{}

func numValidation(n contact) error {
	if n.Name == "" {
		return errors.New("Передано пустое поле")
	} else if len(n.Number) != 11 {
		return errors.New("Передан не валидный номер")
	}
	return nil
}

func handleCreateContact(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()
	var contact contact

	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Ошибка чтения", http.StatusBadRequest)
		return
	}

	if err := numValidation(contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	allContacts = append(allContacts, contact)

	w.WriteHeader(http.StatusCreated)
	b, _ := json.MarshalIndent(contact, "", "    ")
	w.Write(b)
}

func handleAllContacts(w http.ResponseWriter, r *http.Request) {
	defer mtx.RUnlock()
	mtx.RLock()
	if len(allContacts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(allContacts, "", "    ")
	w.Write(b)
}

func handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()
	Name := mux.Vars(r)["Name"]
	var contact contact
	var contactHave bool

	for k, v := range allContacts {
		if Name == v.Name {
			contactHave = true
			if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
				http.Error(w, "Ошибка чтения", http.StatusBadRequest)
				return
			}

			if err := numValidation(contact); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			allContacts[k] = contact
			break
		}

	}
	if !contactHave {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(contact, "", "    ")
	w.Write(b)
}

func handleFavContacts(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()

	var favList []contact
	favoriteFilter := r.URL.Query().Get("favorite")
	if favoriteFilter == "" {
		http.Error(w, "Не указан фильтр поиска", http.StatusBadRequest)
		return
	}

	b, err := strconv.ParseBool(favoriteFilter)
	if err != nil {
		http.Error(w, "Неверно указан фильтр", http.StatusBadRequest)
	}

	for _, v := range allContacts {
		if b == v.IsFav {
			favList = append(favList, v)
		}
	}
	if len(favList) <= 0 {
		http.Error(w, "По указаному фильтру не найдено контактов", http.StatusNotFound)
	}

	f, _ := json.MarshalIndent(favList, "", "    ")
	w.Write(f)
}

func handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	defer mtx.Unlock()
	mtx.Lock()

	Name := mux.Vars(r)["Name"]

	var contactHave bool
	for k, v := range allContacts {
		if Name == v.Name {
			contactHave = true
			allContacts = append(allContacts[:k], allContacts[k+1:]...)
			w.WriteHeader(http.StatusNoContent)
			break
		}
	}
	if !contactHave {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func main() {

	router := mux.NewRouter()

	router.Path("/contacts").Methods("POST").HandlerFunc(handleCreateContact)
	router.Path("/contacts").Methods("GET").HandlerFunc(handleAllContacts)
	router.Path("/contacts/{Name}").Methods("PUT").HandlerFunc(handleUpdateContact)
	router.Path("/contacts/fav").Methods("GET").HandlerFunc(handleFavContacts)
	router.Path("/contacts/{Name}").Methods("DELETE").HandlerFunc(handleDeleteContact)

	if err := http.ListenAndServe(":9091", router); err != nil {
		fmt.Println("Ошибка сервера")
		return
	}
}
