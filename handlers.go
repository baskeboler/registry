package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ServicesInfo{Data: make(map[string]*ServiceRegistry)}
var Services = NewServicesInfo()

func HandleLookup(w http.ResponseWriter, r *http.Request) {
	// r.Context()
	path := r.URL.Path
	serviceName := strings.TrimPrefix(path, "/lookup/")
	registry, ok := Services.Get(serviceName)
	if !ok {
		http.NotFound(w, r)
		return
	}
	header := w.Header()
	header.Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(registry)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var registration = EndpointRegistration{}
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	serviceName := registration.ServiceName
	Services.Set(serviceName, registration)
	w.WriteHeader(http.StatusOK)

}

func HandleDeregister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var registration = EndpointRegistration{}
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	serviceName := registration.ServiceName
	Services.Unset(serviceName, registration)
	w.WriteHeader(http.StatusOK)

}
