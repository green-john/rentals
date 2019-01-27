package http

import "reflect"

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
