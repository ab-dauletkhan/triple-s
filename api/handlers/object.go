package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	objects, err := ReadObjectsFile(bucketName)
	fmt.Println(objects)
	if err != nil {
		log.Printf("Failed to read objects file for bucket %s: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	objectIndex := FindObjectIndex(objects.List, objectKey)
	if objectIndex != -1 {
		objects.List = RemoveObject(objects.List, objectIndex)
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

	newObject := core.Object{
		Name:          objectKey,
		ContentType:   r.Header.Get("Content-Type"),
		ContentLength: strconv.FormatInt(r.ContentLength, 10),
		LastModified:  time.Now().Format(time.RFC3339Nano),
	}
	if newObject.ContentType == "" {
		newObject.ContentType = "application/octet-stream"
	}

	objects.List = append(objects.List, newObject)
	fmt.Println(objects)
	err = WriteObjectsFile(bucketName, objects)
	if err != nil {
		log.Printf("Failed to update objects file for bucket %s: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	log.Printf("Object %s created successfully in bucket %s\n", objectKey, bucketName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object created successfully"))
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
			XMLErrResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	log.Printf("Objects listed successfully for bucket %s\n", bucketName)
	XMLResponse(w, http.StatusOK, objectsData)
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

	objects, err := ReadObjectsFile(bucketName)
	if err != nil {
		log.Printf("Failed to read objects file for bucket %s: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	objectIndex := FindObjectIndex(objects.List, objectKey)
	if objectIndex == -1 {
		log.Printf("Object %s not found in bucket %s\n", objectKey, bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Object not found")
		return
	}

	objects.List = RemoveObject(objects.List, objectIndex)
	err = WriteObjectsFile(bucketName, objects)
	if err != nil {
		log.Printf("Failed to update objects file for bucket %s: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	err = os.Remove(objectPath)
	if err != nil {
		log.Printf("Failed to delete object %s in bucket %s: %v\n", objectKey, bucketName, err)
		XMLErrResponse(w, http.StatusInternalServerError, "Failed to delete object")
		return
	}

	log.Printf("Object %s deleted successfully from bucket %s\n", objectKey, bucketName)
	w.WriteHeader(http.StatusNoContent)
}
