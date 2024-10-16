package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

func CreateObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parsePath(r.URL.Path)
	// if bucketName == "" || objectKey == "" {
	// 	http.Error(w, "Invalid bucket or object key", http.StatusBadRequest)
	// 	return
	// }

	bucketPath := filepath.Join(core.Dir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectPath)
	if err != nil {
		http.Error(w, "Failed to create object", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Failed to write object data", http.StatusInternalServerError)
		return
	}

	err = updateMetadata(bucketPath, objectKey, r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parsePath(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		http.Error(w, "Invalid bucket or object key", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join("data", bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Open(objectPath)
	if err != nil {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	contentType, err := getContentType(bucketPath, objectKey)
	if err != nil {
		http.Error(w, "Failed to get content type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to read object data", http.StatusInternalServerError)
		return
	}
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parsePath(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		http.Error(w, "Invalid bucket or object key", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join("data", bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}

	err := os.Remove(objectPath)
	if err != nil {
		http.Error(w, "Failed to delete object", http.StatusInternalServerError)
		return
	}

	err = removeMetadata(bucketPath, objectKey)
	if err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
