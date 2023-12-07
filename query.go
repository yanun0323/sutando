package sutando

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type query struct {
	coll    *mongo.Collection
	filters []filter
	one     bool
}

func newQuery(collection *mongo.Collection) querying {
	return &query{
		coll:    collection,
		filters: []filter{},
	}
}

func (q *query) build() bson.D {
	query := make(bson.D, 0, len(q.filters))
	for i := range q.filters {
		e := q.filters[i].bson().(bson.E)
		query = append(query, e)
	}
	return query
}

func (q *query) col() *mongo.Collection {
	return q.coll
}

func (q *query) Exists(key string, exists bool) querying {
	return q.appendFilters("$exists", key, exists)
}

func (q *query) And(key string, value any) querying {
	return q.appendFilters("", key, value)
}

func (q *query) Equal(key string, value any) querying {
	return q.appendFilters("$eq", key, value)
}

func (q *query) NotEqual(key string, value ...any) querying {
	return q.appendFilters("$ne", key, value...).Exists(key, true)
}

func (q *query) Greater(key string, value any) querying {
	return q.appendFilters("$gt", key, value)
}

func (q *query) GreaterOrEqual(key string, value any) querying {
	return q.appendFilters("$gte", key, value)
}

func (q *query) Less(key string, value any) querying {
	return q.appendFilters("$lt", key, value)
}

func (q *query) LessOrEqual(key string, value any) querying {
	return q.appendFilters("$lte", key, value)
}

func (q *query) Bitwise(key string, value any) querying {
	return q.appendFilters("$bitsAllSet", key, value)
}

func (q *query) Contain(key string, value ...any) querying {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$all", key, value...)
}

func (q *query) In(key string, value ...any) querying {
	switch len(value) {
	case 0:
		return q
	case 1:
		rv := reflect.ValueOf(value[0])
		if rv.Kind() != reflect.Slice {
			return q.appendFilters("$in", key, value)
		}
		data := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			data = append(data, rv.Index(i).Interface())
		}
		return q.In(key, data...)
	default:
		return q.appendFilters("$in", key, value...)
	}
}

func (q *query) NotIn(key string, value ...any) querying {
	switch len(value) {
	case 0:
		return q
	case 1:
		rv := reflect.ValueOf(value[0])
		if rv.Kind() != reflect.Slice {
			return q.appendFilters("$nin", key, value)
		}
		data := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			data = append(data, rv.Index(i).Interface())
		}
		return q.NotIn(key, data...)
	default:
		return q.appendFilters("$nin", key, value...)
	}
}

func (q *query) Regex(key string, regex string) querying {
	q.appendFilters("$regex", key, regex)
	return q
}

func (q *query) First() querying {
	q.one = true
	return q
}

func (q *query) isOne() bool {
	return q.one
}

func (q *query) appendFilters(operation, key string, value ...any) querying {
	q.filters = append(q.filters, newFilter(operation, key, value...))
	return q
}
