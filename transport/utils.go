package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
)

func contains(a []string, b string) bool {
	for _, elt := range a {
		if elt == b {
			return true
		}
	}

	return false
}

func getJsonTag(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return field.Tag.Get("json")
}

// Utility function to respond to http requests.
func respond(w http.ResponseWriter, status int, data interface{}) {
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
