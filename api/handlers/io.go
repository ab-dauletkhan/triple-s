package handlers

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"github.com/ab-dauletkhan/triple-s/api/core"
)

// ReadBucketsFile reads the buckets meta-file and returns the bucket data
func ReadBucketsFile() (core.Buckets, error) {
	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	log.Printf("Reading buckets meta-file: %s\n", bucketsFilePath)

	records, err := readCSVFile(bucketsFilePath)
	if err != nil {
		return core.Buckets{}, err
	}

	return convertRecordsToBuckets(records), nil
}

// WriteBucketsFile writes the bucket data to the buckets meta-file
func WriteBucketsFile(bucketsData core.Buckets) error {
	bucketsFilePath := filepath.Join(core.Dir, core.BucketsFile)
	records := convertBucketsToRecords(bucketsData)

	return writeCSVFile(bucketsFilePath, core.BucketsCSVHeader, records)
}

// ReadObjectsFile reads the objects file for a bucket and returns the object data
func ReadObjectsFile(bucketName string) (core.Objects, error) {
	objectsFilePath := filepath.Join(core.Dir, bucketName, core.ObjectsFile)
	log.Printf("Reading objects file: %s\n", objectsFilePath)

	records, err := readCSVFile(objectsFilePath)
	if err != nil {
		return core.Objects{}, err
	}

	return convertRecordsToObjects(records), nil
}

// WriteObjectsFile writes the object data to the objects file for a bucket
func WriteObjectsFile(bucketName string, objectsData core.Objects) error {
	objectsFilePath := filepath.Join(core.Dir, bucketName, core.ObjectsFile)
	records := convertObjectsToRecords(objectsData)

	return writeCSVFile(objectsFilePath, core.ObjectsCSVHeader, records)
}

// Reads a CSV file and returns the records without the header
func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		return records[1:], nil // Skip header
	}
	return [][]string{}, nil
}

// Writes a CSV file with the given header and records
func writeCSVFile(filePath string, header []string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if header != nil {
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	if err := writer.WriteAll(records); err != nil {
		return err
	}

	return nil
}
