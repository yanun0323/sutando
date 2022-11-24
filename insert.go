package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insert struct {
	col  *mongo.Collection
	data []bson.M
}

func newInsert(collection *mongo.Collection, p ...any) *insert {
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
	// TODO: [Yanun] ?
	return nil
}

func (ins *insert) optionMany() []*options.InsertManyOptions {
	// TODO: [Yanun] ?
	return nil
}
