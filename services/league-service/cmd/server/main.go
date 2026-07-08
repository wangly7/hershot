package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wangly7/hershot/services/league-service/config"
)

func main() {
	cfg := config.Load()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("league-service ok"))
	})

	port := fmt.Sprintf(":%d", cfg.HTTPPort)

	log.Println("league-service running on ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
