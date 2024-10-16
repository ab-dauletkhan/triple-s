package handlers

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

// XMLErrResponse sends an XML-encoded error response
func XMLErrResponse(w http.ResponseWriter, code int, message string) {
	XMLResponse(w, code, core.Error{Code: code, Message: message})
}

// XMLResponse sends an XML-encoded response
func XMLResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding XML response: %v\n", err)
	}
}

// ReadBucketsFile reads the buckets meta-file and returns the bucket data
func ReadBucketsFile() (core.Buckets, error) {
	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	log.Printf("Reading buckets meta-file: %s\n", bucketsFilePath)

	records, err := readCSVFile(bucketsFilePath)
	if err != nil {
		return core.Buckets{}, fmt.Errorf("error reading buckets file: %w", err)
	}

	return convertRecordsToBuckets(records), nil
}

// WriteBucketsFile writes the bucket data to the buckets meta-file
func WriteBucketsFile(bucketsData core.Buckets) error {
	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	records := convertBucketsToRecords(bucketsData)

	return writeCSVFile(bucketsFilePath, core.BucketsCSVHeader, records)
}

// ReadObjectsFile reads the objects.xml file for a given bucket
func ReadObjectsFile(bucketName string) (core.Objects, error) {
	objectsFilePath := filepath.Join(core.Dir, bucketName, "objects.xml")
	log.Printf("Reading objects.xml file: %s\n", objectsFilePath)

	objectsXML, err := os.ReadFile(objectsFilePath)
	if err != nil {
		return core.Objects{}, fmt.Errorf("error reading objects.xml: %w", err)
	}

	var objectsData core.Objects
	if err := xml.Unmarshal(objectsXML, &objectsData); err != nil {
		return core.Objects{}, fmt.Errorf("error unmarshaling XML from objects.xml: %w", err)
	}

	return objectsData, nil
}

// ParsePath splits the URL path into bucket name and object key
func ParsePath(path string) (bucketName, objectKey string) {
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// UpdateMetadata adds or updates metadata for an object
func UpdateMetadata(bucketPath, objectKey, contentType string) error {
	metadataPath := filepath.Join(bucketPath, "objects.csv")
	return appendCSVRecord(metadataPath, []string{objectKey, contentType})
}

// GetContentType retrieves the content type for an object
func GetContentType(bucketPath, objectKey string) (string, error) {
	metadataPath := filepath.Join(bucketPath, "objects.csv")
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

// RemoveMetadata removes metadata for an object
func RemoveMetadata(bucketPath, objectKey string) error {
	metadataPath := filepath.Join(bucketPath, core.ObjectsFile)
	records, err := readCSVFile(metadataPath)
	if err != nil {
		return fmt.Errorf("error reading metadata file: %w", err)
	}

	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != objectKey {
			updatedRecords = append(updatedRecords, record)
		}
	}

	return writeCSVFile(metadataPath, core.ObjectsCSVHeader, updatedRecords)
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

// CheckBucketEmpty checks if a bucket is empty
func CheckBucketEmpty(bucketName string) error {
	objectsData, err := ReadObjectsFile(bucketName)
	if err != nil {
		return fmt.Errorf("error reading objects file: %w", err)
	}

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

// CreateBucketDirectory creates a directory for a bucket
func CreateBucketDirectory(bucketName string) error {
	bucketDirPath := filepath.Join(core.Dir, bucketName)
	return os.MkdirAll(bucketDirPath, core.DirPerm)
}

// HandleError handles different types of errors and sends appropriate XML responses
func HandleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrBucketNotFound:
		XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
	case ErrBucketAlreadyExists:
		XMLErrResponse(w, http.StatusConflict, "Bucket already exists")
	case ErrBucketNotEmpty:
		XMLErrResponse(w, http.StatusConflict, "Bucket is not empty")
	default:
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Helper functions for CSV operations

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV from %s: %w", filePath, err)
	}

	if len(records) > 1 {
		return records[1:], nil // Skip header
	}
	return [][]string{}, nil
}

func writeCSVFile(filePath string, header []string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if header != nil {
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing header to %s: %w", filePath, err)
		}
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("error writing records to %s: %w", filePath, err)
	}

	return nil
}

func appendCSVRecord(filePath string, record []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("error appending record to %s: %w", filePath, err)
	}

	return nil
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
