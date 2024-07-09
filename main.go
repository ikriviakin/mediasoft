package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	// Определение маршрутов
	router.HandleFunc("/{objects}", CreateEntity).Methods("POST")
	router.HandleFunc("/{objects}", UpdateEntities).Methods("PUT")
	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	router.HandleFunc("/cars/{id}", UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", PatchCar).Methods("PATCH")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	// Запуск сервера
	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
