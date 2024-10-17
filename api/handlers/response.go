package handlers

import (
	"encoding/xml"
	"log"
	"net/http"

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
