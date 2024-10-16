package types

import (
	"encoding/xml"
	"time"
)

type Bucket struct {
	XMLName    xml.Name  `xml:"Bucket"`
	Name       string    `xml:"Name"`
	CreatedAt  time.Time `xml:"CreatedAt"`
	ModifiedAt time.Time `xml:"ModifiedAt"`
	Status     string    `xml:"Status"`
}

type Buckets struct {
	XMLName xml.Name `xml:"Buckets"`
	Buckets []Bucket `xml:"Bucket"`
}
