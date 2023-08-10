package sutando

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (ins *insert) optionOne() []*options.InsertOneOptions {
	// TODO: Implement me
	return nil
}

func (ins *insert) optionMany() []*options.InsertManyOptions {
	// TODO: Implement me
	return nil
}

func (ins *insert) Exec(ctx context.Context) (insertOneResult, insertManyResult, error) {
	var (
		err  error
		one  insertOneResult
		many insertManyResult
	)
	objects := ins.build()
	switch len(objects) {
	case 0:
		return one, many, errors.New("object to insert should be pointer")
	case 1:
		one, err = ins.col.InsertOne(ctx, objects[0], ins.optionOne()...)
		return one, many, err
	default:
		many, err = ins.col.InsertMany(ctx, objects, ins.optionMany()...)
		return one, many, err
	}
}
