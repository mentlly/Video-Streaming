package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("GO_VIDEO_MANAGEMENT_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "Ok"}`)
	})
	log.Printf("Go Video Management Server booting up internally on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start Go server: %v", err)
	}
}
