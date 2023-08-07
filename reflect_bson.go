package sutando

import (
	"reflect"
	"unsafe"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

func bsonEncoder(p any, key string) bson.M {
	var elem reflect.Value

	switch reflect.TypeOf(p).Kind() {
	case reflect.Struct:
		elem = reflect.ValueOf(p)
		if elem.Kind() != reflect.Struct {
			return nil
		}
	case reflect.Pointer:
		elem = reflect.ValueOf(p).Elem()
		if elem.Kind() != reflect.Struct {
			return nil
		}
	}

	result := make(bson.M, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)

		label, skip, inline, omitempty := getTag(elem, field)

		if skip {
			continue
		}

		fValue := elem.Field(i)
		if omitempty && fValue.IsZero() {
			continue
		}

		if inline {
			result[label] = bsonEncoder(fValue.Interface(), "")
			continue
		}

		result[label] = valueParse(elem, i)
	}

	if len(key) > 0 {
		return bson.M{key: result}
	}

	return result
}

func valueParse(elem reflect.Value, fieldIndex int) any {
	v := elem.Field(fieldIndex)
	if elem.Type().Field(fieldIndex).IsExported() {
		if v.Type() == _TYPE_DECIMAL {
			return v.Interface().(decimal.Decimal).String()
		}
		return v.Interface()
	}
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface()
}
