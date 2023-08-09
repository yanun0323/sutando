package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (d *delete) Bitwise(key string, value any) deleting {
	d.q.Bitwise(key, value)
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

func (d *delete) First() deleting {
	d.q.First()
	return d
}

func (d *delete) Exec(ctx context.Context) (deleteResult, error) {
	if d.q.isOne() {
		return d.q.col().DeleteOne(ctx, d.q.build(), &options.DeleteOptions{})
	}
	return d.q.col().DeleteMany(ctx, d.q.build(), &options.DeleteOptions{})
}
