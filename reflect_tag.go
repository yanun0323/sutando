package sutando

import (
	"reflect"
	"strings"
)

const (
	_tagKey       string = "bson"
	_tagOmitempty string = "omitempty"
	_tagIgnore    string = "-"
	_tagID        string = "_id"
)

func getTag(field reflect.StructField) (label string, skip bool, omitempty bool) {
	tags := strings.Split(field.Tag.Get(_tagKey), ",")
	label = tags[0]
	if label == _tagIgnore || label == _tagID {
		return "", true, false
	}

	if len(label) == 0 {
		label = firstLowerCase(field.Name)
	}

	if len(tags) == 1 {
		return label, false, false
	}

	for _, tag := range tags[1:] {
		omitempty = omitempty || tag == _tagOmitempty
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
