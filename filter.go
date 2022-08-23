package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	_FILTER_TYPE = reflect.TypeOf(filter{})
)

type filter struct {
	operation string
	key       string
	values    []any
}

func newFilter(operation, key string, v ...any) filter {
	if len(v) > 0 {
		return filter{operation, key, v}
	}
	return filter{operation, key, nil}
}

func (f filter) bson() any {
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

func (f filter) bsonKeyOperation() any {
	if len(f.values) > 1 || f.operation == "$all" {
		return bson.E{Key: f.key, Value: bson.M{f.operation: decodeValues(f.values)}}
	}
	if len(f.values) == 1 {
		return bson.E{Key: f.key, Value: bson.M{f.operation: decodeValues(f.values)[0]}}
	}
	return nil
}

func (f filter) bsonOnlyOperation() any {
	if len(f.values) > 1 || f.operation == "$all" {
		return bson.M{f.operation: decodeValues(f.values)}
	}
	if len(f.values) == 1 {
		return bson.M{f.operation: decodeValues(f.values)[0]}
	}
	return nil
}

func (f filter) bsonOnlyKey() any {
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
			result = append(result, a.(filter).bson())
			continue
		}
		result = append(result, a)
	}
	return result
}
