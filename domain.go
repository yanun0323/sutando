package sutando

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type common interface {
	col() *mongo.Collection
	build() bson.D
	isOne() bool
}

type query interface {
	common

	Exists(key string, exists bool) query
	And(key string, value any) query
	Equal(key string, value any) query
	NotEqual(key string, value ...any) query
	Greater(key string, value any) query
	GreaterOrEqual(key string, value any) query
	Less(key string, value any) query
	LessOrEqual(key string, value any) query
	Bitwise(key string, value any) query
	Contain(key string, value ...any) query
	In(key string, value ...any) query
	NotIn(key string, value ...any) query
	First() query
}

type update interface {
	common

	Exists(key string, exists bool) update
	And(key string, value any) update
	Equal(key string, value any) update
	NotEqual(key string, value ...any) update
	Greater(key string, value any) update
	GreaterOrEqual(key string, value any) update
	Less(key string, value any) update
	LessOrEqual(key string, value any) update
	Bitwise(key string, value any) update
	Contain(key string, value ...any) update
	In(key string, value ...any) update
	NotIn(key string, value ...any) update
	First() update

	Set(key string, value any) update

	buildObjects() any
}
