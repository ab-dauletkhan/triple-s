package handlers

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

// ================================================================================================
func XMLErrResponse(w http.ResponseWriter, code int, message string) {
	XMLResponse(w, code, core.Error{Code: code, Message: message})
}

func XMLResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding XML response: %v\n", err)
	}
}

// func XMLResponse(w http.ResponseWriter, code int, data interface{}) {
// 	xmlResp, err := xml.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/xml")
// 	w.WriteHeader(code)
// 	w.Write(xmlResp)
// }

// ================================================================================================
func ReadBucketsFile() (core.Buckets, error) {
	var bucketsData core.Buckets

	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	log.Printf("Reading buckets meta-file: %s\n", bucketsFilePath)
	csvFile, err := os.Open(bucketsFilePath)
	if err != nil {
		return bucketsData, fmt.Errorf("error opening %s: %w", bucketsFilePath, err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return bucketsData, fmt.Errorf("error reading CSV records from %s: %w", bucketsFilePath, err)
	}

	if len(records) > 1 {
		records = records[1:]
	}

	for _, record := range records {
		bucket := core.Bucket{
			Name:         record[0],
			Status:       record[1],
			CreationDate: record[2],
			LastUpdated:  record[3],
		}
		bucketsData.List = append(bucketsData.List, bucket)
	}

	return bucketsData, nil
}

func WriteBucketsFile(bucketsData core.Buckets) error {
	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	file, err := os.Create(bucketsFilePath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", bucketsFilePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	records, err := convertBucketsToRecords(bucketsData)
	if err != nil {
		return fmt.Errorf("error converting buckets data to records: %w", err)
	}

	err = writer.Write(core.BucketsCSVHeader)
	if err != nil {
		return fmt.Errorf("error writing header to %s: %w", bucketsFilePath, err)
	}

	err = writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", bucketsFilePath, err)
	}
	return nil
}

func convertBucketsToRecords(bucketsData core.Buckets) ([][]string, error) {
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
	return records, nil
}

// ================================================================================================
func ReadObjectsFile(bucketName string) (core.Objects, error) {
	var objectsData core.Objects
	log.Printf("Reading objects.xml file: %s\n", core.Dir+"/"+bucketName+"/objects.xml")
	objectsXML, err := os.ReadFile(core.Dir + "/" + bucketName + "/objects.xml")
	if err != nil {
		return objectsData, fmt.Errorf("error reading objects.xml: %w", err)
	}
	err = xml.Unmarshal(objectsXML, &objectsData)
	if err != nil {
		return objectsData, fmt.Errorf("error unmarshaling XML from objects.xml: %w", err)
	}
	return objectsData, nil
}

func parsePath(path string) (string, string) {
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func updateMetadata(bucketPath, objectKey, contentType string) error {
	metadataPath := filepath.Join(bucketPath, "objects.csv")
	file, err := os.OpenFile(metadataPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write([]string{objectKey, contentType})
}

func getContentType(bucketPath, objectKey string) (string, error) {
	metadataPath := filepath.Join(bucketPath, "objects.csv")
	file, err := os.Open(metadataPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if record[0] == objectKey {
			return record[1], nil
		}
	}
	return "", errors.New("content type not found")
}

func removeMetadata(bucketPath, objectKey string) error {
	metadataPath := filepath.Join(bucketPath, "objects.csv")
	file, err := os.Open(metadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var records [][]string
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if record[0] != objectKey {
			records = append(records, record)
		}
	}

	file, err = os.Create(metadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(records)
}

// ================================================================================================
func parseW3CTime(timestamp string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func findBucketIndex(buckets []core.Bucket, name string) int {
	for i, bucket := range buckets {
		if bucket.Name == name {
			return i
		}
	}
	return -1
}

func checkBucketEmpty(bucketName string) error {
	objectsData, err := ReadObjectsFile(bucketName)
	if err != nil {
		return fmt.Errorf("error reading objects file: %w", err)
	}

	if len(objectsData.List) > 0 {
		return ErrBucketNotEmpty
	}

	return nil
}

func removeBucket(buckets []core.Bucket, index int) []core.Bucket {
	return append(buckets[:index], buckets[index+1:]...)
}

func createBucketDirectory(bucketName string) error {
	bucketDirPath := filepath.Join(core.Dir, bucketName)
	if _, err := os.Stat(bucketDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(bucketDirPath, core.DirPerm); err != nil {
			return fmt.Errorf("error creating bucket directory %s: %w", bucketDirPath, err)
		}
	}
	return nil
}

func handleError(w http.ResponseWriter, err error) {
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
