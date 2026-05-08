package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/undndnwnkk/yadro-test-task/internal/dns"
)

var service = dns.NewService("/etc/resolv.conf")

func main() {
	http.HandleFunc("/dns", dnsHandler)
	log.Println("Server starting on :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dnsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST: %s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		list, err := service.GetServers()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var request struct {
			IP string `json:"ip"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request", 400)
			return
		}

		if err := service.AddServer(request.IP); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete:
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "missing ip param", 400)
			return
		}

		if err := service.RemoveServer(ip); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

}
