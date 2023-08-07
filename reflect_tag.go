package sutando

import (
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	_TAG_KEY          string = "bson"
	_TAG_INLINE       string = "inline"
	_TAG_OMITEMPTY    string = "omitempty"
	_TAG_IGNORE_VALUE string = "-"
)

var (
	_TYPE_DECIMAL = reflect.TypeOf(decimal.Decimal{})
)

func getTag(v reflect.Value, field reflect.StructField) (label string, skip bool, inline bool, omitempty bool) {
	tags := strings.Split(field.Tag.Get(_TAG_KEY), ",")
	label = tags[0]
	if label == _TAG_IGNORE_VALUE {
		return "", true, false, false
	}

	if len(label) == 0 {
		label = firstLowerCase(field.Name)
	}

	if len(tags) == 1 {
		return label, false, false, false
	}

	for _, tag := range tags[1:] {
		inline = inline || tag == _TAG_INLINE
		omitempty = omitempty || tag == _TAG_OMITEMPTY
	}

	return label, false, inline, omitempty
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
