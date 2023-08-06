package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type common interface {
	col() *mongo.Collection
	build() bson.D
	isOne() bool
}

type inserting interface {
	Exec(context.Context) (insertOneResult, insertManyResult, error)
}

type updating interface {
	Exists(key string, exists bool) updating
	And(key string, value any) updating
	Equal(key string, value any) updating
	NotEqual(key string, value ...any) updating
	Greater(key string, value any) updating
	GreaterOrEqual(key string, value any) updating
	Less(key string, value any) updating
	LessOrEqual(key string, value any) updating
	Bitwise(key string, value any) updating
	Contain(key string, value ...any) updating
	In(key string, value ...any) updating
	NotIn(key string, value ...any) updating
	First() updating

	Set(key string, value any) updating
	Exec(ctx context.Context, upsert bool) (updateResult, error)

	buildObjects() any
}

type querying interface {
	common

	Exists(key string, exists bool) querying
	And(key string, value any) querying
	Equal(key string, value any) querying
	NotEqual(key string, value ...any) querying
	Greater(key string, value any) querying
	GreaterOrEqual(key string, value any) querying
	Less(key string, value any) querying
	LessOrEqual(key string, value any) querying
	Bitwise(key string, value any) querying
	Contain(key string, value ...any) querying
	In(key string, value ...any) querying
	NotIn(key string, value ...any) querying
	First() querying
}

type finding interface {
	Exists(key string, exists bool) finding
	And(key string, value any) finding
	Equal(key string, value any) finding
	NotEqual(key string, value ...any) finding
	Greater(key string, value any) finding
	GreaterOrEqual(key string, value any) finding
	Less(key string, value any) finding
	LessOrEqual(key string, value any) finding
	Bitwise(key string, value any) finding
	Contain(key string, value ...any) finding
	In(key string, value ...any) finding
	NotIn(key string, value ...any) finding
	First() finding

	Exec(ctx context.Context, result any) error
}

type deleting interface {
	Exists(key string, exists bool) deleting
	And(key string, value any) deleting
	Equal(key string, value any) deleting
	NotEqual(key string, value ...any) deleting
	Greater(key string, value any) deleting
	GreaterOrEqual(key string, value any) deleting
	Less(key string, value any) deleting
	LessOrEqual(key string, value any) deleting
	Bitwise(key string, value any) deleting
	Contain(key string, value ...any) deleting
	In(key string, value ...any) deleting
	NotIn(key string, value ...any) deleting
	First() deleting

	Exec(context.Context) (deleteResult, error)
}
