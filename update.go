package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type updateQ struct {
	data []bson.M
	q    query
}

func newUpdate(collection *mongo.Collection, p ...any) update {
	d := make([]bson.M, 0, len(p))
	for i := range p {
		if reflect.TypeOf(p[i]).Kind() != reflect.Pointer {
			continue
		}
		d = append(d, bsonEncoder(p[i], reflect.TypeOf(p).Name(), true))
	}
	return &updateQ{
		q:    newFind(collection),
		data: d,
	}
}

func (u *updateQ) col() *mongo.Collection {
	return u.q.col()
}

func (u *updateQ) build() bson.D {
	return u.q.build()
}

func (u *updateQ) buildObjects() []any {
	result := make([]any, 0, len(u.data))
	for i := range u.data {
		result = append(result, u.data[i])
	}
	return result
}

func (u *updateQ) Exists(key string, exists bool) update {
	u.q.Exists(key, exists)
	return u
}

func (u *updateQ) And(key string, value any) update {
	u.q.And(key, value)
	return u
}

func (u *updateQ) Equal(key string, value any) update {
	u.q.Equal(key, value)
	return u
}

func (u *updateQ) NotEqual(key string, value ...any) update {
	u.q.NotEqual(key, value...)
	return u
}

func (u *updateQ) Greater(key string, value any) update {
	u.q.Greater(key, value)
	return u
}

func (u *updateQ) GreaterOrEqual(key string, value any) update {
	u.q.GreaterOrEqual(key, value)
	return u
}

func (u *updateQ) Less(key string, value any) update {
	u.q.Less(key, value)
	return u
}

func (u *updateQ) LessOrEqual(key string, value any) update {
	u.q.LessOrEqual(key, value)
	return u
}

func (u *updateQ) Bitwise(key string, value any) update {
	u.q.Bitwise(key, value)
	return u
}

func (u *updateQ) Contain(key string, value ...any) update {
	u.q.Contain(key, value...)
	return u
}

func (u *updateQ) In(key string, value ...any) update {
	u.q.In(key, value...)
	return u
}

func (u *updateQ) NotIn(key string, value ...any) update {
	u.q.NotIn(key, value...)
	return u
}
