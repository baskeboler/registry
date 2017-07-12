package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// var _ = sync.Map

func main() {
	var (
		port string
		ok   bool
	)
	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}
	router := mux.NewRouter()
	router.HandleFunc("/lookup/{service}", HandleLookup).Methods("GET")
	router.HandleFunc("/register", HandleRegister).Methods("POST")
	router.HandleFunc("/deregister", HandleDeregister).Methods("POST")
	handler := handlers.LoggingHandler(os.Stdout, router)
	log.Printf("Starting Registry on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
