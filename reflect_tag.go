package sutando

import (
	"reflect"

	"github.com/shopspring/decimal"
)

const (
	_TAG_KEY          string = "bson"
	_TAG_IGNORE_VALUE string = "-"
)

var (
	_TYPE_DECIMAL = reflect.TypeOf(decimal.Decimal{})
)

func getTag(v reflect.Value, field reflect.StructField) (label string, skip bool) {
	label = field.Tag.Get(_TAG_KEY)
	if label == _TAG_IGNORE_VALUE {
		return "", true
	}

	if label == "" {
		label = firstLowerCase(field.Name)
	}
	return label, false
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
