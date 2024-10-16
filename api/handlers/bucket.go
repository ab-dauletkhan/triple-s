package handlers

import (
	"fmt"
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
	ErrBucketNotFound      = fmt.Errorf("bucket not found")
	ErrBucketAlreadyExists = fmt.Errorf("bucket already exists")
	ErrBucketNotEmpty      = fmt.Errorf("bucket is not empty")
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	newBucket, err := createBucket(bucketName)
	if err != nil {
		log.Printf("Error creating bucket %s: %v\n", bucketName, err)
		handleError(w, err)
		return
	}

	log.Printf("Bucket %s created successfully", bucketName)
	XMLResponse(w, http.StatusOK, newBucket)
}

func createBucket(bucketName string) (core.Bucket, error) {
	if err := util.ValidateBucketName(bucketName); err != nil {
		return core.Bucket{}, fmt.Errorf("invalid bucket name: %w", err)
	}

	bucketsData, err := ReadBucketsFile()
	if err != nil {
		return core.Bucket{}, fmt.Errorf("error reading buckets file: %w", err)
	}

	if findBucketIndex(bucketsData.List, bucketName) != -1 {
		return core.Bucket{}, ErrBucketAlreadyExists
	}

	newBucket := core.Bucket{
		Name:         bucketName,
		Status:       "Active",
		CreationDate: time.Now().Format(time.RFC3339Nano),
		LastUpdated:  time.Now().Format(time.RFC3339Nano),
	}
	bucketsData.List = append(bucketsData.List, newBucket)

	if err := WriteBucketsFile(bucketsData); err != nil {
		return core.Bucket{}, fmt.Errorf("error writing buckets file: %w", err)
	}

	if err := createBucketDirectory(bucketName); err != nil {
		return core.Bucket{}, err
	}

	if err := util.InitObjectFile(bucketName); err != nil {
		return core.Bucket{}, fmt.Errorf("error initializing object file: %w", err)
	}

	return newBucket, nil
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	bucketsData, err := listBuckets()
	if err != nil {
		log.Printf("Error listing buckets: %v\n", err)
		handleError(w, err)
		return
	}

	log.Println("Buckets listed successfully")
	XMLResponse(w, http.StatusOK, bucketsData)
}

func listBuckets() (core.Buckets, error) {
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		return core.Buckets{}, fmt.Errorf("error reading buckets file: %w", err)
	}
	return bucketsData, nil
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	if err := deleteBucket(bucketName); err != nil {
		log.Printf("Error deleting bucket %s: %v\n", bucketName, err)
		handleError(w, err)
		return
	}

	log.Printf("Bucket %s deleted successfully", bucketName)
	w.WriteHeader(http.StatusNoContent)
}

func deleteBucket(bucketName string) error {
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		return fmt.Errorf("error reading buckets file: %w", err)
	}

	bucketIndex := findBucketIndex(bucketsData.List, bucketName)
	if bucketIndex == -1 {
		return ErrBucketNotFound
	}

	if err := checkBucketEmpty(bucketName); err != nil {
		return err
	}

	bucketsData.List = removeBucket(bucketsData.List, bucketIndex)

	if err := WriteBucketsFile(bucketsData); err != nil {
		return fmt.Errorf("error writing buckets file: %w", err)
	}

	bucketDirPath := filepath.Join(core.Dir, bucketName)
	if err := os.RemoveAll(bucketDirPath); err != nil {
		return fmt.Errorf("error deleting bucket directory %s: %w", bucketDirPath, err)
	}

	return nil
}
