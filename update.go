package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type updateQ struct {
	data bson.M
	q    query
	one  bool
	set  bson.M
}

func newUpdate(collection *mongo.Collection, p any) update {
	var d bson.M = nil
	if p != nil {
		d = bsonEncoder(p, reflect.TypeOf(p).Name())
	}
	return &updateQ{
		q:    newFind(collection),
		data: d,
		set:  bson.M{},
	}
}

func (u *updateQ) col() *mongo.Collection {
	return u.q.col()
}

func (u *updateQ) build() bson.D {
	return u.q.build()
}

func (u *updateQ) buildObjects() any {
	if u.data != nil {
		return bson.M{"$set": u.data}
	}
	return bson.M{"$set": u.set}
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

func (u *updateQ) First() update {
	u.one = true
	return u
}

func (u *updateQ) isOne() bool {
	return u.one
}

func (u *updateQ) Set(key string, value any) update {
	u.set[key] = value
	return u
}
