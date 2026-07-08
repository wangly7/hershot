package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8083"

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ingestion-service ok"))
	})

	log.Println("ingestion-service running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
