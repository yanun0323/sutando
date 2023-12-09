package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	_FILTER_TYPE = reflect.TypeOf(filters{})
)

type filters struct {
	operation string
	key       string
	values    []any
}

func newFilters(operation, key string, v ...any) filters {
	if len(v) > 0 {
		return filters{operation, key, v}
	}
	return filters{operation, key, nil}
}

func (f filters) bson() any {
	if len(f.operation) > 0 && len(f.key) > 0 {
		return f.bsonKeyOperation()
	}

	if len(f.operation) > 0 && len(f.key) == 0 {
		return f.bsonOnlyOperation()
	}

	if len(f.operation) == 0 && len(f.key) > 0 {
		return f.bsonOnlyKey()
	}
	return nil
}

func (f filters) bsonKeyOperation() any {
	if len(f.values) > 1 || f.operation == "$all" {
		return bson.E{Key: f.key, Value: bson.M{f.operation: decodeValues(f.values)}}
	}
	if len(f.values) == 1 {
		return bson.E{Key: f.key, Value: bson.M{f.operation: decodeValues(f.values)[0]}}
	}
	return nil
}

func (f filters) bsonOnlyOperation() any {
	if len(f.values) > 1 || f.operation == "$all" {
		return bson.M{f.operation: decodeValues(f.values)}
	}
	if len(f.values) == 1 {
		return bson.M{f.operation: decodeValues(f.values)[0]}
	}
	return nil
}

func (f filters) bsonOnlyKey() any {
	if len(f.values) == 1 {
		return bson.E{Key: f.key, Value: decodeValues(f.values)[0]}
	}
	if len(f.values) > 1 {
		return bson.E{Key: f.key, Value: decodeValues(f.values)}
	}
	return nil
}

func decodeValues(values []any) []any {
	result := make([]any, 0, len(values))
	for _, a := range values {
		if a == nil {
			continue
		}
		if reflect.TypeOf(a) == _FILTER_TYPE {
			result = append(result, a.(filters).bson())
			continue
		}
		result = append(result, a)
	}
	return result
}
