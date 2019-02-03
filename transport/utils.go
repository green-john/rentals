package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)


// Utility function to respond to http requests.
func respond(w http.ResponseWriter, status int, data interface{}) {
	if p, ok := data.(Public); ok {
		data = p.Public()
	}

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := io.Copy(w, &buffer); err != nil {
		log.Println("Error responding:", err)
	}
}
