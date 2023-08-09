package sutando

import (
	"reflect"
	"strings"
)

const (
	_TAG_KEY       string = "bson"
	_TAG_OMITEMPTY string = "omitempty"
	_TAG_IGNORE    string = "-"
	_TAG_ID        string = "_id"
)

func getTag(v reflect.Value, field reflect.StructField) (label string, skip bool, omitempty bool) {
	tags := strings.Split(field.Tag.Get(_TAG_KEY), ",")
	label = tags[0]
	if label == _TAG_IGNORE || label == _TAG_ID {
		return "", true, false
	}

	if len(label) == 0 {
		label = firstLowerCase(field.Name)
	}

	if len(tags) == 1 {
		return label, false, false
	}

	for _, tag := range tags[1:] {
		omitempty = omitempty || tag == _TAG_OMITEMPTY
	}

	return label, false, omitempty
}

func firstLowerCase(s string) string {
	if s == "" {
		return s
	}

	str := []byte(s)
	upper := str[0] >= 'A' && str[0] <= 'Z'
	if !upper {
		return s
	}

	str[0] += 32
	return string(str)
}
