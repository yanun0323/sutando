package sutando

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _zeroID = primitive.ObjectID{}

type insertResult struct {
	InsertedIDs []any
}

func newInsertedResult(count int, a ...any) insertResult {
	ids := make([]any, 0, count)
	ids = append(ids, a...)
	if offset := count - len(ids); offset > 0 {
		for i := 0; i < offset; i++ {
			ids = append(ids, _zeroID)
		}
	}
	return insertResult{
		InsertedIDs: ids,
	}
}

type insert struct {
	col     *mongo.Collection
	data    []bson.M
	encoder bsonEncoder
}

func newInsert(collection *mongo.Collection, encoder bsonEncoder, p ...any) inserting {
	ins := &insert{
		col:     collection,
		data:    make([]bson.M, 0, len(p)),
		encoder: encoder,
	}

	name := reflect.TypeOf(p).Name()
	for i := range p {
		ins.handleData(p[i], name)
	}

	return ins
}

func (ins *insert) handleData(elem any, name string) {
	rValue := reflect.ValueOf(elem)
	rType := reflect.TypeOf(elem)
	switch rValue.Kind() {
	case reflect.Slice:
		rName := rType.Name()
		for i := 0; i < rValue.Len(); i++ {
			ins.handleData(rValue.Index(i).Interface(), rName)
		}
	case reflect.Pointer:
		ins.handleData(rValue.Elem().Interface(), name)
	case reflect.Struct:
		ins.data = append(ins.data, ins.encoder.Encode(elem, name))
	case reflect.Map, reflect.Func:
		return
	default:
		return
	}
}

func (ins *insert) build() []any {
	result := make([]any, 0, len(ins.data))
	for i := range ins.data {
		result = append(result, ins.data[i])
	}
	return result
}

func (ins *insert) Exec(ctx context.Context) (insertResult, error) {
	objects := ins.build()
	objCount := len(objects)
	switch objCount {
	case 0:
		return insertResult{}, errors.New("mongo: object to insert should be pointer")
	case 1:
		one, err := ins.col.InsertOne(ctx, objects[0])
		return newInsertedResult(objCount, one.InsertedID), ins.wrapDuplicateKeyErr(err)
	default:
		many, err := ins.col.InsertMany(ctx, objects)
		return newInsertedResult(objCount, many.InsertedIDs...), ins.wrapDuplicateKeyErr(err)
	}
}

func (ins *insert) wrapDuplicateKeyErr(err error) error {
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicatedKey
	}
	return err
}
