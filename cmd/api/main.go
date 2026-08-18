package main

import (
	"contacts-api/internal/contacts/feature/service"
	postgres "contacts-api/internal/contacts/feature/storage"
	"contacts-api/internal/contacts/feature/transport"
	"contacts-api/internal/core/logger"
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	var ctx = context.Background()

	log, logFileClose, err := logger.NewLoger("DEBUG")
	if err != nil {
		log.Fatal("error to create log")
		return
	}
	defer logFileClose()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	conn, err := postgres.CreateConnection(ctx)
	if err != nil {
		log.Error("failed to connection db")
		return
	}
	log.Info("database connect succesfully")

	storage := postgres.NewStorage(conn)
	service := service.NewContactService(storage, log)
	handler := transport.NewHandler(service, log)

	router := mux.NewRouter()

	router.Path("/contacts").Methods("POST").HandlerFunc(handler.CreateContact)
	router.Path("/contacts").Methods("GET").HandlerFunc(handler.AllContacts)
	router.Path("/contacts/{ID}").Methods("PUT").HandlerFunc(handler.UpdateContact)
	router.Path("/contacts/fav").Methods("GET").HandlerFunc(handler.FavContacts)
	router.Path("/contacts/{ID}").Methods("DELETE").HandlerFunc(handler.DeleteContact)

	go func() {
		if err := http.ListenAndServe(":9091", router); err != nil {
			fmt.Println("Ошибка сервера", err)
			return
		}
	}()

	<-ctx.Done()
	log.Info("App stop correctly")
}
