package sutando

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type find struct {
	col     *mongo.Collection
	filters []filter
}

func newFind(collection *mongo.Collection) *find {
	return &find{
		col:     collection,
		filters: []filter{}}
}

func (q *find) build() bson.D {
	find := make(bson.D, 0, len(q.filters))
	for i := range q.filters {
		e := q.filters[i].bson().(bson.E)
		find = append(find, e)
	}
	return find
}

func (q *find) Exists(key string, exists bool) *find {
	return q.appendFilters("$exists", key, exists)
}

func (q *find) And(key string, value any) *find {
	return q.appendFilters("", key, value)
}

func (q *find) Equal(key string, value any) *find {
	return q.appendFilters("$eq", key, value)
}

func (q *find) NotEqual(key string, value ...any) *find {
	return q.appendFilters("$ne", key, value...).Exists(key, true)
}

func (q *find) Greater(key string, value any) *find {
	return q.appendFilters("$gt", key, value)
}

func (q *find) GreaterOrEqual(key string, value any) *find {
	return q.appendFilters("$gte", key, value)
}

func (q *find) Less(key string, value any) *find {
	return q.appendFilters("$lt", key, value)
}

func (q *find) LessOrEqual(key string, value any) *find {
	return q.appendFilters("$lte", key, value)
}

func (q *find) Bitwise(key string, value any) *find {
	return q.appendFilters("$bitsAllSet", key, value)
}

func (q *find) Contain(key string, value ...any) *find {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$all", key, value...)
}

func (q *find) In(key string, value ...any) *find {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$in", key, value...)
}

func (q *find) NotIn(key string, value ...any) *find {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$nin", key, value...)
}

func (q *find) appendFilters(operation, key string, value ...any) *find {
	q.filters = append(q.filters, newFilter(operation, key, value...))
	return q
}
