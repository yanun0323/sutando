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
	col  *mongo.Collection
	data []bson.M
}

func newInsert(collection *mongo.Collection, p ...any) inserting {
	d := make([]bson.M, 0, len(p))
	for i := range p {
		if reflect.TypeOf(p[i]).Kind() != reflect.Pointer {
			continue
		}
		d = append(d, bsonEncoder(p[i], reflect.TypeOf(p).Name()))
	}
	return &insert{
		col:  collection,
		data: d,
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
