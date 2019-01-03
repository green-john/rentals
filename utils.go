package rentals

func contains(a []string, b string) bool {
	for _, elt := range a {
		if elt == b {
			return true
		}
	}

	return false
}
