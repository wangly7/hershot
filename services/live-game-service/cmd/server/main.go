package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8082"

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("live-game-service ok"))
	})

	log.Println("live-game-service running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
