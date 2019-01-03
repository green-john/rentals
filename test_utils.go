package rentals

import "testing"

func ok(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assert(t *testing.T, expr bool, errorMsg string) {
	if !expr {
		t.Errorf(errorMsg)
	}
}
