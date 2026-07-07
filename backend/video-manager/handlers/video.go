package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"video-manager/services"
	"video-manager/utils"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	services.VideoProccessor(dir)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}
