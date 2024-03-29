package sutando

import (
	"context"

	"github.com/yanun0323/sutando/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deleteResult *mongo.DeleteResult

type delete struct {
	q querying
}

func newDelete(collection *mongo.Collection) deleting {
	return &delete{
		q: newQuery(collection),
	}
}

func (d *delete) Exists(key string, exists bool) deleting {
	d.q.Exists(key, exists)
	return d
}

func (d *delete) And(key string, value any) deleting {
	d.q.And(key, value)
	return d
}

func (d *delete) Equal(key string, value any) deleting {
	d.q.Equal(key, value)
	return d
}

func (d *delete) NotEqual(key string, value ...any) deleting {
	d.q.NotEqual(key, value...)
	return d
}

func (d *delete) Greater(key string, value any) deleting {
	d.q.Greater(key, value)
	return d
}

func (d *delete) GreaterOrEqual(key string, value any) deleting {
	d.q.GreaterOrEqual(key, value)
	return d
}

func (d *delete) Less(key string, value any) deleting {
	d.q.Less(key, value)
	return d
}

func (d *delete) LessOrEqual(key string, value any) deleting {
	d.q.LessOrEqual(key, value)
	return d
}

func (d *delete) Contain(key string, value ...any) deleting {
	d.q.Contain(key, value...)
	return d
}

func (d *delete) In(key string, value ...any) deleting {
	d.q.In(key, value...)
	return d
}

func (d *delete) NotIn(key string, value ...any) deleting {
	d.q.NotIn(key, value...)
	return d
}

func (d *delete) Regex(key string, regex string, opt ...option.Regex) deleting {
	d.q.Regex(key, regex, opt...)
	return d
}

func (d *delete) Bson(e ...bson.E) deleting {
	d.q.Bson(e...)
	return d
}

func (d *delete) First() deleting {
	d.q.first()
	return d
}

func (d *delete) Exec(ctx context.Context) (deleteResult, error) {
	if d.q.isOne() {
		return d.q.col().DeleteOne(ctx, d.q.build(), &options.DeleteOptions{})
	}
	return d.q.col().DeleteMany(ctx, d.q.build(), &options.DeleteOptions{})
}
