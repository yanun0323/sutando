package sutando

import (
	"context"
	"reflect"

	"github.com/yanun0323/sutando/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type query struct {
	coll *mongo.Collection
	f    filter
	one  bool
}

func newQuery(collection *mongo.Collection) querying {
	return &query{
		coll: collection,
		f:    filter{d: bson.D{}},
	}
}

func (q *query) build() bson.D {
	return q.f.d
}

func (q *query) col() *mongo.Collection {
	return q.coll
}

func (q *query) Exists(key string, exists bool) querying {
	return q.add(key, bson.M{"$exists": exists})
}

func (q *query) And(key string, value any) querying {
	return q.add(key, value)
}

func (q *query) Equal(key string, value any) querying {
	return q.add(key, bson.M{"$eq": value})
}

func (q *query) NotEqual(key string, value ...any) querying {
	if len(value) == 0 {
		return q
	}
	return q.add(key, bson.M{"$ne": purge(value)}).Exists(key, true)
}

func (q *query) Greater(key string, value any) querying {
	return q.add(key, bson.M{"$gt": value})
}

func (q *query) GreaterOrEqual(key string, value any) querying {
	return q.add(key, bson.M{"$gte": value})
}

func (q *query) Less(key string, value any) querying {
	return q.add(key, bson.M{"$lt": value})
}

func (q *query) LessOrEqual(key string, value any) querying {
	return q.add(key, bson.M{"$lte": value})
}

func (q *query) Bitwise(key string, value any) querying {
	return q.add(key, bson.M{"$bitsAllSet": value})
}

func (q *query) Contain(key string, value ...any) querying {
	if len(value) == 0 {
		return q
	}
	return q.add(key, bson.M{"$all": value})
}

func (q *query) In(key string, value ...any) querying {
	switch len(value) {
	case 0:
		return q
	case 1:
		rv := reflect.ValueOf(value[0])
		if rv.Kind() != reflect.Slice {
			return q.add(key, bson.M{"$in": value})
		}
		data := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			data = append(data, rv.Index(i).Interface())
		}
		return q.In(key, data...)
	default:
		return q.add(key, bson.M{"$in": value})
	}
}

func (q *query) NotIn(key string, value ...any) querying {
	switch len(value) {
	case 0:
		return q
	case 1:
		rv := reflect.ValueOf(value[0])
		if rv.Kind() != reflect.Slice {
			return q.add(key, bson.M{"$nin": value})
		}
		data := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			data = append(data, rv.Index(i).Interface())
		}
		return q.NotIn(key, data...)
	default:
		return q.add(key, bson.M{"$nin": value})
	}
}

func (q *query) Regex(key string, regex string, opt ...option.Regex) querying {
	if len(opt) == 0 || opt[0] == 0 {
		return q.add(key, bson.M{"$regex": regex})
	}
	o := opt[0]
	buf := make([]byte, 0, 4)
	if o&option.CaseInsensitive == 1 {
		buf = append(buf, 'i')
	}

	if o&option.MatchMultiLine == 1 {
		buf = append(buf, 'm')
	}

	if o&option.IgnoreWhitespace == 1 {
		buf = append(buf, 'x')
	}

	if o&option.DotMatchNewLine == 1 {
		buf = append(buf, 's')
	}

	return q.add(key, bson.M{"$regex": regex, "$options": string(buf)})
}

func (q *query) Bson(e ...bson.E) querying {
	q.f.appends(e...)
	return q
}

func (q *query) Count(ctx context.Context, index ...string) (int64, error) {
	opt := options.Count()
	if len(index) != 0 && len(index[0]) != 0 {
		opt.SetHint(index[0])
	}
	i, err := q.col().CountDocuments(ctx, q.build(), opt)
	if err != nil {
		return q.col().CountDocuments(ctx, q.build(), options.Count())
	}
	return i, nil
}

func (q *query) first() querying {
	q.one = true
	return q
}

func (q *query) isOne() bool {
	return q.one
}

func (q *query) add(key string, val any) *query {
	q.f.append(key, val)
	return q
}
