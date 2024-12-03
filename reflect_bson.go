package sutando

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	_TYPE_TIME    = reflect.TypeOf(time.Time{})
	_TYPE_DECIMAL = reflect.TypeOf(decimal.Decimal{})
)

type bsonEncoder struct {
	supportedEncodeTypes map[reflect.Type]bool
}

func newBsonEncoder(supportedEncodeTypes ...reflect.Type) bsonEncoder {
	m := make(map[reflect.Type]bool, len(supportedEncodeTypes))
	for _, t := range supportedEncodeTypes {
		m[t] = true
	}
	return bsonEncoder{
		supportedEncodeTypes: m,
	}
}

func (e bsonEncoder) Encode(p any, key string) bson.M {
	switch reflect.TypeOf(p).Kind() {
	case reflect.Pointer:
		return e.Encode(reflect.ValueOf(p).Elem().Interface(), key)
	case reflect.Struct:
		return e.encode(p, key).(bson.M)
	case reflect.Map:
		return p.(bson.M)
	default:
		return bson.M{}
	}
}

func (e bsonEncoder) encode(p any, key string) any {
	switch reflect.TypeOf(p).Kind() {
	case reflect.Slice, reflect.Array:
		rv := reflect.ValueOf(p)
		res := make(bson.A, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			res = append(res, e.encode(rv.Index(i).Interface(), ""))
		}
		return e.resultWrapper(res, key)
	case reflect.Map:
		rv := reflect.ValueOf(p)
		res := make(bson.M, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			res[k.String()] = e.encode(v.Interface(), "")
		}
		return e.resultWrapper(res, key)
	case reflect.Pointer:
		return e.encode(reflect.ValueOf(p).Elem().Interface(), key)
	case reflect.Struct:
		elem := reflect.ValueOf(p)

		result := make(bson.M, elem.NumField())
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)

			label, skip, omitempty := getTag(field)
			if skip {
				continue
			}

			fValue := elem.Field(i)
			if omitempty && fValue.IsZero() {
				continue
			}

			if e.isSupportedStructure(fValue.Interface()) {
				result[label] = e.valueParse(elem, i)
				continue
			}

			result[label] = e.encode(fValue.Interface(), "")
		}

		return e.resultWrapper(result, key)
	}

	return reflect.ValueOf(p).Interface()
}

func (e bsonEncoder) resultWrapper(result any, key string) any {
	if len(key) > 0 {
		return bson.M{key: result}
	}
	return result
}

func (e bsonEncoder) valueParse(elem reflect.Value, fieldIndex int) any {
	v := elem.Field(fieldIndex)
	if elem.Type().Field(fieldIndex).IsExported() {
		if v.Type() == _TYPE_DECIMAL {
			return v.Interface().(decimal.Decimal).String()
		}
		return v.Interface()
	}
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface()
}

func (e bsonEncoder) isSupportedStructure(p any) bool {
	rt := reflect.TypeOf(p)
	rk := rt.Kind()

	switch rk {
	case reflect.Pointer:
		return e.isSupportedStructure(reflect.ValueOf(p).Elem().Interface())
	case reflect.Struct:
		switch rt {
		case _TYPE_DECIMAL, _TYPE_TIME:
			return true
		default:
			return e.supportedEncodeTypes[rt]
		}
	case reflect.Array, reflect.Slice, reflect.Map:
		return false
	default:
		return true
	}
}
