package sutando

import (
	"context"
	"reflect"

	"github.com/yanun0323/sutando/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type update struct {
	data bson.M
	q    querying
	set  bson.M
}

func newUpdate(collection *mongo.Collection, encoder bsonEncoder, p any) updating {
	var d bson.M = nil
	if p != nil {
		d = encoder.Encode(p, reflect.TypeOf(p).Name())
	}
	return &update{
		q:    newQuery(collection),
		data: d,
		set:  bson.M{},
	}
}

func (u *update) buildObjects() any {
	if u.data != nil {
		return bson.M{"$set": u.data}
	}
	return bson.M{"$set": u.set}
}

func (u *update) Exists(key string, exists bool) updating {
	u.q.Exists(key, exists)
	return u
}

func (u *update) And(key string, value any) updating {
	u.q.And(key, value)
	return u
}

func (u *update) Equal(key string, value any) updating {
	u.q.Equal(key, value)
	return u
}

func (u *update) NotEqual(key string, value ...any) updating {
	u.q.NotEqual(key, value...)
	return u
}

func (u *update) Greater(key string, value any) updating {
	u.q.Greater(key, value)
	return u
}

func (u *update) GreaterOrEqual(key string, value any) updating {
	u.q.GreaterOrEqual(key, value)
	return u
}

func (u *update) Less(key string, value any) updating {
	u.q.Less(key, value)
	return u
}

func (u *update) LessOrEqual(key string, value any) updating {
	u.q.LessOrEqual(key, value)
	return u
}

func (u *update) Bitwise(key string, value any) updating {
	u.q.Bitwise(key, value)
	return u
}

func (u *update) Contain(key string, value ...any) updating {
	u.q.Contain(key, value...)
	return u
}

func (u *update) In(key string, value ...any) updating {
	u.q.In(key, value...)
	return u
}

func (u *update) NotIn(key string, value ...any) updating {
	u.q.NotIn(key, value...)
	return u
}

func (u *update) Regex(key string, regex string, opt ...option.Regex) updating {
	u.q.Regex(key, regex, opt...)
	return u
}

func (u *update) First() updating {
	u.q.First()
	return u
}

func (u *update) Set(key string, value any) updating {
	u.set[key] = value
	return u
}
func (u *update) Exec(ctx context.Context, upsert bool) (updateResult, error) {
	if u.q.isOne() {
		return u.q.col().UpdateOne(ctx, u.q.build(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	}
	return u.q.col().UpdateMany(ctx, u.q.build(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})

}
