package rentals

import "testing"

func ok(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func assert(t *testing.T, expr bool, errorMsg string) {
	t.Helper()

	if !expr {
		t.Errorf(errorMsg)
	}
}
