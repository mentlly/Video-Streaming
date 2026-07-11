package main

import (
	"log"
	"mime"
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

	mime.AddExtensionType(".m3u8", "application/x-mpegURL")
	mime.AddExtensionType(".ts", "video/MP2T")
	hlsDir := http.Dir("./uploads")
	fileServer := http.FileServer(hlsDir)
	mux.Handle("GET /api/video/", http.StripPrefix("/api/video/", fileServer))

	log.Printf("Go Video Management Server booting up internally on port %s", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Failed to start Go server: %v", err)
	}
}
