package sutando

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type query struct {
	col     *mongo.Collection
	filters []filter
}

func newQuery(collection *mongo.Collection) *query {
	return &query{
		col:     collection,
		filters: []filter{}}
}

func (q *query) build() bson.D {
	query := make(bson.D, 0, len(q.filters))
	for i := range q.filters {
		e := q.filters[i].bson().(bson.E)
		query = append(query, e)
	}
	return query
}

func (q *query) Exists(key string, exists bool) *query {
	return q.appendFilters("$exists", key, exists)
}

func (q *query) And(key string, value any) *query {
	return q.appendFilters("", key, value)
}

func (q *query) Equal(key string, value any) *query {
	return q.appendFilters("$eq", key, value)
}

func (q *query) NotEqual(key string, value ...any) *query {
	return q.appendFilters("$ne", key, value...).Exists(key, true)
}

func (q *query) Greater(key string, value any) *query {
	return q.appendFilters("$gt", key, value)
}

func (q *query) GreaterOrEqual(key string, value any) *query {
	return q.appendFilters("$gte", key, value)
}

func (q *query) Less(key string, value any) *query {
	return q.appendFilters("$lt", key, value)
}

func (q *query) LessOrEqual(key string, value any) *query {
	return q.appendFilters("$lte", key, value)
}

func (q *query) Bitwise(key string, value any) *query {
	return q.appendFilters("$bitsAllSet", key, value)
}

func (q *query) Contain(key string, value ...any) *query {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$all", key, value...)
}

func (q *query) In(key string, value ...any) *query {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$in", key, value...)
}

func (q *query) NotIn(key string, value ...any) *query {
	if len(value) == 0 {
		return q
	}
	return q.appendFilters("$nin", key, value...)
}

func (q *query) appendFilters(operation, key string, value ...any) *query {
	q.filters = append(q.filters, newFilter(operation, key, value...))
	return q
}
