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

	"github.com/ab-dauletkhan/triple-s/api/core"
	"github.com/ab-dauletkhan/triple-s/api/types"
)

func XMLErrResponse(w http.ResponseWriter, code int, message string) {
	XMLResponse(w, code, types.Error{Code: code, Message: message})
}

func XMLResponse(w http.ResponseWriter, code int, data interface{}) {
	xmlResp, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write(xmlResp)
}

func ReadBucketsFile() (types.Buckets, error) {
	var bucketsData types.Buckets
	log.Printf("Reading buckets.xml file: %s\n", core.Dir+"/buckets.xml")
	bucketsXML, err := os.ReadFile(core.Dir + "/buckets.xml")
	if err != nil {
		return bucketsData, fmt.Errorf("error reading buckets.xml: %w", err)
	}
	err = xml.Unmarshal(bucketsXML, &bucketsData)
	if err != nil {
		return bucketsData, fmt.Errorf("error unmarshaling XML from buckets.xml: %w", err)
	}
	return bucketsData, nil
}

func WriteBucketsFile(bucketsData types.Buckets) error {
	updatedXML, err := xml.MarshalIndent(bucketsData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling updated XML: %w", err)
	}
	err = os.WriteFile(core.Dir+"/buckets.xml", updatedXML, 0o644)
	if err != nil {
		return fmt.Errorf("error writing buckets.xml: %w", err)
	}
	return nil
}

func ReadObjectsFile(bucketName string) (types.Objects, error) {
	var objectsData types.Objects
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
	file, err := os.OpenFile(metadataPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
