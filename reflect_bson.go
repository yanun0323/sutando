package sutando

import (
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

func bsonEncoder(p any, key string, update bool) bson.M {
	elem := reflect.ValueOf(p).Elem()
	if elem.Kind() != reflect.Struct {
		return nil
	}

	result := make(bson.M, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)

		label, skip := getTag(elem, field)
		if skip {
			continue
		}
		result[label] = valueParse(elem.Field(i))
	}

	if update {
		result = bson.M{"$set": result}
	}

	if len(key) > 0 {
		return bson.M{key: result}
	}
	return result
}

func valueParse(v reflect.Value) any {
	if v.Type() == _TYPE_DECIMAL {
		return v.Interface().(decimal.Decimal).String()
	}
	return v.Interface()
}
