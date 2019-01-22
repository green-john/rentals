package e2e

import (
	"encoding/json"
	"fmt"
	"rentals/tst"
	"testing"
)

func loginWithUser(t *testing.T, serverUrl, username, pwd string) (string, error) {
	t.Helper()

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, pwd)
	response, err := tst.MakeRequest("POST", serverUrl+"/login", "", []byte(body))
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 399 {
		t.Errorf("Got %d code, expected 2XX", response.StatusCode)
	}

	var token struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&token)

	if err != nil {
		return "", err
	}

	return token.Token, nil
}
