package validation

import (
	"reflect"
	"strings"
)

// FieldName returns the JSON tag name if present; otherwise the struct field name.
func FieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	name := strings.Split(jsonTag, ",")[0]
	if name == "" {
		name = field.Name
	}
	return name
}

// ValidateByTag scans struct fields and calls check(field, value) for each field whose validate tag contains tagContains.
// If check returns true, the field is appended to invalid list.
func ValidateByTag(v interface{}, tagContains string, check func(field reflect.StructField, value reflect.Value) bool) []string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	var invalid []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fv := val.Field(i)
		validateTag := field.Tag.Get("validate")
		if !strings.Contains(validateTag, tagContains) {
			continue
		}
		if check(field, fv) {
			invalid = append(invalid, FieldName(field))
		}
	}
	return invalid
}
