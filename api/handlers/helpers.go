package handlers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

// ParsePath splits the URL path into bucket name and object key
func ParsePath(path string) (bucketName, objectKey string) {
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// FindBucketIndex finds the index of a bucket in a slice of buckets
func FindBucketIndex(buckets []core.Bucket, name string) int {
	for i, bucket := range buckets {
		if bucket.Name == name {
			return i
		}
	}
	return -1
}

// FindObjectIndex finds the index of an object in a slice of objects
func FindObjectIndex(objects []core.Object, name string) int {
	for i, object := range objects {
		if object.Name == name {
			return i
		}
	}
	return -1
}

// CheckBucketEmpty checks if a bucket is empty
func CheckBucketEmpty(bucketName string) error {
	objectsRecords, err := readCSVFile(filepath.Join(core.Dir, bucketName, core.ObjectsFile))
	if err != nil {
		return fmt.Errorf("error reading objects file: %w", err)
	}

	objectsData := convertRecordsToObjects(objectsRecords)
	if len(objectsData.List) > 0 {
		return ErrBucketNotEmpty
	}

	return nil
}

// RemoveBucket removes a bucket from a slice of buckets
func RemoveBucket(buckets []core.Bucket, index int) []core.Bucket {
	buckets[index] = buckets[len(buckets)-1]
	return buckets[:len(buckets)-1]
}

// RemoveObject removes an object from a slice of objects
func RemoveObject(objects []core.Object, index int) []core.Object {
	objects[index] = objects[len(objects)-1]
	return objects[:len(objects)-1]
}

// CreateBucketDirectory creates a directory for a bucket
func CreateBucketDirectory(bucketName string) error {
	bucketDirPath := filepath.Join(core.Dir, bucketName)
	return os.MkdirAll(bucketDirPath, core.DirPerm)
}

// GetContentType retrieves the content type for an object
func GetContentType(bucketPath, objectKey string) (string, error) {
	metadataPath := filepath.Join(bucketPath, core.ObjectsFile)
	records, err := readCSVFile(metadataPath)
	if err != nil {
		return "", fmt.Errorf("error reading metadata file: %w", err)
	}

	for _, record := range records {
		if record[0] == objectKey {
			return record[1], nil
		}
	}

	return "", errors.New("content type not found")
}

func convertRecordsToBuckets(records [][]string) core.Buckets {
	var bucketsData core.Buckets
	for _, record := range records {
		bucket := core.Bucket{
			Name:         record[0],
			Status:       record[1],
			CreationDate: record[2],
			LastUpdated:  record[3],
		}
		bucketsData.List = append(bucketsData.List, bucket)
	}
	return bucketsData
}

func convertBucketsToRecords(bucketsData core.Buckets) [][]string {
	var records [][]string
	for _, bucket := range bucketsData.List {
		record := []string{
			bucket.Name,
			bucket.Status,
			bucket.CreationDate,
			bucket.LastUpdated,
		}
		records = append(records, record)
	}
	return records
}

func convertObjectsToRecords(objectsData core.Objects) [][]string {
	var records [][]string
	for _, object := range objectsData.List {
		record := []string{
			object.Name,
			object.ContentType,
			object.ContentLength,
			object.LastModified,
		}
		records = append(records, record)
	}
	return records
}

func convertRecordsToObjects(records [][]string) core.Objects {
	var objectsData core.Objects
	for _, record := range records {
		object := core.Object{
			Name:          record[0],
			ContentType:   record[1],
			ContentLength: record[2],
			LastModified:  record[3],
		}
		objectsData.List = append(objectsData.List, object)
	}
	return objectsData
}
