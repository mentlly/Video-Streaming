package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"video-manager/services"
	"video-manager/utils"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
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

	title := r.FormValue("title")
	if strings.TrimSpace(title) != "" {
		title = handler.Filename
	}
	description := r.FormValue("description")
	channel_id := r.FormValue("channel_id")

	//Creating directory for uploaded video to write in
	videoId := utils.GenerateVideoId()
	dir := "./uploads/" + videoId
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Internal server error creating object folder", http.StatusBadRequest)
	}

	dst, err := os.Create(filepath.Join(dir, "original.mp4"))
	if err != nil {
		http.Error(w, "Internal server error creating file", http.StatusBadRequest)
		return
	}
	defer dst.Close()

	//Copying the file to the created filename
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusBadRequest)
	}

	duration := services.VideoProccessor(dir)
	services.UploadVideoDb(account_id, channel_id, videoId, title, description, duration)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}
