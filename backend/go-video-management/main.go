package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("GO_VIDEO_MANAGEMENT_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "Ok"}`)
	})

	//Creating upload folder to store uploaded videos
	os.MkdirAll("./uploads", os.ModePerm)

	http.HandleFunc("/api/upload", uploadVideoHandler)

	log.Printf("Go Video Management Server booting up internally on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start Go server: %v", err)
	}
}

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

func uploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "File to large or bad request", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	fmt.Printf("Uploaded File: %s\n", handler.Filename)
	fmt.Printf("File Size: %d bytes\n", handler.Size)

	dst, err := os.Create(filepath.Join("./uploads", handler.Filename))
	if err != nil {
		http.Error(w, "Internal server error creating file", http.StatusBadRequest)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}
