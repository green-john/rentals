package transport

import "net/http"

func ErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errorMsg))
}
