package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/undndnwnkk/yadro-test-task/internal/dns"
)

var service = dns.NewService("/etc/resolv.conf")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dns", dnsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server starting on :8080...")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	log.Println("Press Ctrl+C to stop")

	<-done
	log.Println("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Println("Server exited properly")
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
