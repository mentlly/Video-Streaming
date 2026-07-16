package main

import (
	"log"
	"net/http"
	"os"

	"video-manager/handlers"
	"video-manager/services"
)

func main() {
	port := os.Getenv("GO_VIDEO_MANAGEMENT_PORT")
	if port == "" {
		port = "8080"
	}

	services.InitDb()

	//Creating upload folder to store uploaded videos
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		log.Fatalf("Failed to create base upload directory: %v", err)
	}

	mux := http.NewServeMux()

	protectedRoute := services.AuthMiddleware(http.HandlerFunc(handlers.UploadVideoHandler))
	mux.Handle("POST /api/video/upload", protectedRoute)

	mux.HandleFunc("GET /api/video/get/latest", handlers.GetLatestVideoHandler)
	mux.HandleFunc("GET /api/video/get/{VideoId}/{filename...}", handlers.GetIdVideoHandler)

	log.Printf("Go Video Management Server booting up internally on port %s", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Failed to start Go server: %v", err)
	}
}
