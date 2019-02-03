package tst

import (
	"bytes"
	"net/http"
	"testing"
)

// Package tst has common utilities for testing

func Ok(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func True(t *testing.T, expr bool, errorMsg string) {
	t.Helper()

	if !expr {
		t.Errorf(errorMsg)
	}
}

func MakeRequest(method, url, authToken string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set json as default
	req.Header.Set("Content-Type", "application/json")

	// If an auth token is provided, add it to the headers
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	return http.DefaultClient.Do(req)
}
