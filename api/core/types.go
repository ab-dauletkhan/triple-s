package core

import (
	"encoding/xml"
)

type Bucket struct {
	XMLName      xml.Name `xml:"Bucket"`
	Name         string   `xml:"Name"`
	CreationDate string   `xml:"CreationDate"`
	LastUpdated  string   `xml:"LastUpdated"`
	Status       string   `xml:"Status"`
}

type Buckets struct {
	XMLName xml.Name `xml:"Buckets"`
	List    []Bucket `xml:"Bucket"`
}

type Object struct {
	XMLName       xml.Name `xml:"Object"`
	Name          string   `xml:"Name"`
	ContentType   string   `xml:"ContentType"`
	ContentLength string   `xml:"ContentLength"`
}

type Objects struct {
	XMLName xml.Name `xml:"Objects"`
	List    []Object `xml:"Object"`
}

type Error struct {
	Code    int    `xml:"Code"`
	Message string `xml:"Message"`
	// TODO: add resource field
}
