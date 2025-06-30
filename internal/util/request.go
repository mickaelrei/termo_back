package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// ReadBody reads the request body and unmarshals it into the given result. If anything goes wrong, it writes an error
// to the response writer. Returns whether it succeeded
//
// This function is supposed to be used in a handler function like this:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    var result struct { ... }
//	    if !util.ReadBody(w, r, &result) {
//	        return
//	    }
//	    ...
//	}
func ReadBody[T any](w http.ResponseWriter, r *http.Request, result *T) bool {
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[io.ReadAll] | %v", err)
		WriteInternalError(w)
		return false
	}

	// Unmarshal body
	err = json.Unmarshal(body, result)
	if err != nil {
		// Check if it's a syntax error
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			http.Error(w, "Invalid JSON syntax", http.StatusBadRequest)
			return false
		}

		// Check if it's a decoding error
		var decodingErr *json.UnmarshalTypeError
		if errors.As(err, &decodingErr) {
			http.Error(w, "Invalid JSON type", http.StatusBadRequest)
			return false
		}

		// Unknown error
		WriteInternalError(w)
		return false
	}

	return true
}

// WriteResponseJSON writes a JSON response. If anything goes wrong, it writes an error to the response writer.
func WriteResponseJSON[T any](w http.ResponseWriter, response T) {
	// Convert structure into string bytes
	bytes, err := json.Marshal(response)
	if err != nil {
		WriteInternalError(w)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Write response
	_, err = w.Write(bytes)
	if err != nil {
		WriteInternalError(w)
		return
	}
}

// WriteInternalError writes a default message for when an internal error occurred
func WriteInternalError(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}
