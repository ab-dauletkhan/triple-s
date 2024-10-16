package handlers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ab-dauletkhan/triple-s/api/core"
	"github.com/ab-dauletkhan/triple-s/api/types"
	"github.com/ab-dauletkhan/triple-s/api/util"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	log.Println("Create bucket called")

	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	// 1. Validate the bucket name
	err := util.ValidateBucketName(bucketName)
	if err != nil {
		log.Printf("Invalid bucket name: %s, error: %s\n", bucketName, err.Error())
		XMLErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid bucket name: %s", err.Error()))
		return
	}

	// 2. Read the buckets.xml file
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 3. Check if the bucket already exists in the buckets.xml file
	for _, bucket := range bucketsData.Buckets {
		if bucket.Name == bucketName {
			log.Printf("Bucket %s already exists\n", bucketName)
			XMLErrResponse(w, http.StatusConflict, "Bucket already exists")
			return
		}
	}

	// 4. Create the new bucket entry
	newBucket := types.Bucket{
		Name:       bucketName,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Status:     "Active",
	}
	bucketsData.Buckets = append(bucketsData.Buckets, newBucket)

	// 5. Write the updated buckets.xml file
	err = WriteBucketsFile(bucketsData)
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 6. Create a directory for the new bucket (if necessary)
	bucketDirPath := core.Dir + "/" + bucketName
	if _, err := os.Stat(bucketDirPath); os.IsNotExist(err) {
		err = os.Mkdir(bucketDirPath, 0o755)
		if err != nil {
			log.Printf("error creating bucket directory %s: %s\n", bucketDirPath, err.Error())
			XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	// 7. Create an objects.xml file for the new bucket
	err = util.InitObjectFile(bucketName)
	if err != nil {
		log.Printf("error creating objects.xml for bucket %s: %s\n", bucketName, err.Error())
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 8. Respond with success
	log.Printf("Bucket %s created successfully", bucketName)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	XMLResponse(w, http.StatusOK, newBucket)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	log.Println("List buckets called")

	// 1. Read the buckets.xml file
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 2. Marshal the buckets data to XML
	bucketsXML, err := xml.MarshalIndent(bucketsData, "", "  ")
	if err != nil {
		log.Printf("error marshaling XML: %s\n", err.Error())
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	log.Println("Responded with\n", string(bucketsXML))
	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	w.Write(bucketsXML)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	// TODO: Order of the delete is confusing, what happens if some step fails, how to properly handle them
	// 1. Get the bucket name from the URL
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	// 2. Read the buckets.xml file
	bucketsData, err := ReadBucketsFile()
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 3. Check if the bucket exists in the buckets.xml file
	bucketIndex := -1
	for i, bucket := range bucketsData.Buckets {
		if bucket.Name == bucketName {
			bucketIndex = i
			break
		}
	}

	if bucketIndex == -1 {
		log.Printf("Bucket %s not found\n", bucketName)
		XMLErrResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	// 4. Check if the bucket is empty
	objectsData, err := ReadObjectsFile(bucketName)
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(objectsData.Objects) > 0 {
		log.Printf("Bucket %s is not empty\n", bucketName)
		XMLErrResponse(w, http.StatusConflict, "Bucket is not empty")
		return
	}

	// 5. Delete the bucket entry from the buckets
	bucketsData.Buckets[bucketIndex] = bucketsData.Buckets[len(bucketsData.Buckets)-1]
	bucketsData.Buckets = bucketsData.Buckets[:len(bucketsData.Buckets)-1]

	// 6. Write the updated buckets.xml file
	err = WriteBucketsFile(bucketsData)
	if err != nil {
		log.Println(err)
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 7. Delete the bucket directory
	bucketDirPath := core.Dir + "/" + bucketName
	err = os.RemoveAll(bucketDirPath)
	if err != nil {
		log.Printf("error deleting bucket directory %s: %s\n", bucketDirPath, err.Error())
		XMLErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// 8. Respond with success
	log.Printf("Bucket %s deleted successfully", bucketName)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Bucket '%s' deleted successfully", bucketName)))
}
