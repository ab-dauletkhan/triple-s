package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ab-dauletkhan/triple-s/api/core"
	"github.com/ab-dauletkhan/triple-s/api/util"
)

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrBucketNotEmpty      = errors.New("bucket is not empty")
	ErrInternalServer      = errors.New("internal server error")
)

// CreateBucket creates a new bucket
// 1. Extract the bucket name from the URL path
// 2. Validate the bucket name
// 3. Read the buckets file
// 4. Check if the bucket already exists
// 5. Create a new bucket
// 6. Write the updated buckets data to the file
// 7. Create a directory for the new bucket
// 8. Initialize the object file for the new bucket
func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	if err := util.ValidateBucketName(bucketName); err != nil {
		log.Printf("Error validating bucket name %s: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Printf("Error reading buckets file: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	if FindBucketIndex(bucketsData.List, bucketName) != -1 {
		log.Printf("Bucket %s already exists\n", bucketName)
		XMLErrResponse(w, http.StatusConflict, ErrBucketAlreadyExists.Error())
		return
	}

	newBucket := core.Bucket{
		Name:         bucketName,
		Status:       "Active",
		CreationDate: time.Now().Format(time.RFC3339Nano),
		LastUpdated:  time.Now().Format(time.RFC3339Nano),
	}
	bucketsData.List = append(bucketsData.List, newBucket)

	if err := WriteBucketsFile(bucketsData); err != nil {
		log.Printf("Error writing buckets file: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	if err := CreateBucketDirectory(bucketName); err != nil {
		log.Printf("Error creating bucket directory: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	if err := util.InitObjectFile(bucketName); err != nil {
		log.Printf("Error initializing object file: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	log.Printf("Bucket %s created successfully", bucketName)
	XMLResponse(w, http.StatusOK, newBucket)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Printf("error reading buckets file: %s", err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	log.Println("Buckets listed successfully")
	XMLResponse(w, http.StatusOK, bucketsData)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Printf("Error reading buckets file: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	bucketIndex := FindBucketIndex(bucketsData.List, bucketName)
	if bucketIndex == -1 {
		log.Printf("Bucket %s not found\n", bucketName)
		XMLErrResponse(w, http.StatusNotFound, ErrBucketNotFound.Error())
		return
	}

	if err := CheckBucketEmpty(bucketName); err != nil {
		log.Printf("Error checking if bucket %s is empty: %v\n", bucketName, err)
		XMLErrResponse(w, http.StatusConflict, ErrBucketNotEmpty.Error())
		return
	}

	bucketsData.List = RemoveBucket(bucketsData.List, bucketIndex)

	if err := WriteBucketsFile(bucketsData); err != nil {
		log.Printf("Error writing buckets file: %v\n", err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	bucketDirPath := filepath.Join(core.Dir, bucketName)
	if err := os.RemoveAll(bucketDirPath); err != nil {
		log.Printf("Error deleting bucket directory %s: %v\n", bucketDirPath, err)
		XMLErrResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	log.Printf("Bucket %s deleted successfully", bucketName)
	w.WriteHeader(http.StatusNoContent)
}
