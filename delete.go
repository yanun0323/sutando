package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type delete struct {
	col  *mongo.Collection
	data []bson.M
}

func newDelete(collection *mongo.Collection, p ...any) *delete {
	d := make([]bson.M, 0, len(p))
	for i := range p {
		if reflect.TypeOf(p[i]).Kind() != reflect.Pointer {
			continue
		}
		d = append(d, bsonEncoder(p[i], reflect.TypeOf(p).Name(), false))
	}
	return &delete{
		col:  collection,
		data: d,
	}
}
