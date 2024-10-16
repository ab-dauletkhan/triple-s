package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

func CreateObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ParsePath(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		log.Println("Invalid bucket or object key")
		XMLErrResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := filepath.Join(core.Dir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectPath)
	if err != nil {
		log.Printf("Failed to create object %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to create object")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Printf("Failed to write object data for %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to write object data")
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err = UpdateMetadata(bucketPath, objectKey, contentType)
	if err != nil {
		log.Printf("Failed to update metadata for %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to update metadata")
		return
	}

	log.Printf("Object %s created successfully in bucket %s\n", objectKey, bucketName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object created successfully"))
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ParsePath(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		log.Println("Invalid bucket or object key")
		XMLErrResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := filepath.Join(core.Dir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Open(objectPath)
	if err != nil {
		log.Printf("Object not found: %s in bucket %s\n", objectKey, bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Object not found")
		return
	}
	defer file.Close()

	contentType, err := GetContentType(bucketPath, objectKey)
	if err != nil {
		log.Printf("Failed to get content type for %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to get content type")
		return
	}

	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Failed to read object data for %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to read object data")
		return
	}

	log.Printf("Object %s retrieved successfully from bucket %s\n", objectKey, bucketName)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ParsePath(r.URL.Path)
	if bucketName == "" || objectKey == "" {
		log.Println("Invalid bucket or object key")
		XMLErrResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := filepath.Join(core.Dir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		log.Printf("Object not found: %s in bucket %s\n", objectKey, bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Object not found")
		return
	}

	err := os.Remove(objectPath)
	if err != nil {
		log.Printf("Failed to delete object %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to delete object")
		return
	}

	err = RemoveMetadata(bucketPath, objectKey)
	if err != nil {
		log.Printf("Failed to update metadata after deleting %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to update metadata")
		return
	}

	log.Printf("Object %s deleted successfully from bucket %s\n", objectKey, bucketName)
	w.WriteHeader(http.StatusNoContent)
}

func ListObjects(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	objectsData, err := ReadObjectsFile(bucketName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Bucket not found: %s\n", bucketName)
			XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
		} else {
			log.Printf("Error reading objects file for bucket %s: %v\n", bucketName, err)
			XMLErrResponse(w, http.StatusInternalServerError, "Error reading objects file")
		}
		return
	}

	log.Printf("Objects listed successfully for bucket %s\n", bucketName)
	XMLResponse(w, http.StatusOK, objectsData)
}
