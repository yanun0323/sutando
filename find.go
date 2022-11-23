package sutando

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type find struct {
	coll    *mongo.Collection
	filters []filter
}

func newFind(collection *mongo.Collection) query {
	return &find{
		coll:    collection,
		filters: []filter{}}
}

func (f *find) build() bson.D {
	query := make(bson.D, 0, len(f.filters))
	for i := range f.filters {
		e := f.filters[i].bson().(bson.E)
		query = append(query, e)
	}
	return query
}

func (f *find) col() *mongo.Collection {
	return f.coll
}

func (f *find) Exists(key string, exists bool) query {
	return f.appendFilters("$exists", key, exists)
}

func (f *find) And(key string, value any) query {
	return f.appendFilters("", key, value)
}

func (f *find) Equal(key string, value any) query {
	return f.appendFilters("$eq", key, value)
}

func (f *find) NotEqual(key string, value ...any) query {
	return f.appendFilters("$ne", key, value...).Exists(key, true)
}

func (f *find) Greater(key string, value any) query {
	return f.appendFilters("$gt", key, value)
}

func (f *find) GreaterOrEqual(key string, value any) query {
	return f.appendFilters("$gte", key, value)
}

func (f *find) Less(key string, value any) query {
	return f.appendFilters("$lt", key, value)
}

func (f *find) LessOrEqual(key string, value any) query {
	return f.appendFilters("$lte", key, value)
}

func (f *find) Bitwise(key string, value any) query {
	return f.appendFilters("$bitsAllSet", key, value)
}

func (f *find) Contain(key string, value ...any) query {
	if len(value) == 0 {
		return f
	}
	return f.appendFilters("$all", key, value...)
}

func (f *find) In(key string, value ...any) query {
	if len(value) == 0 {
		return f
	}
	return f.appendFilters("$in", key, value...)
}

func (f *find) NotIn(key string, value ...any) query {
	if len(value) == 0 {
		return f
	}
	return f.appendFilters("$nin", key, value...)
}

func (f *find) appendFilters(operation, key string, value ...any) query {
	f.filters = append(f.filters, newFilter(operation, key, value...))
	return f
}
