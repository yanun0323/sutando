package sutando

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type update struct {
	col     *mongo.Collection
	filters []filter
	keys    []string
	data    []bson.M
	id      uint64
}

func newUpdate(collection *mongo.Collection, p ...any) *update {
	d := make([]bson.M, 0, len(p))
	for i := range p {
		fmt.Println(reflect.TypeOf(p[i]).Kind())
		if reflect.TypeOf(p[i]).Kind() != reflect.Pointer {
			continue
		}
		d = append(d, bsonEncoder(p[i], reflect.TypeOf(p).Name(), true))
	}
	return &update{
		col:     collection,
		filters: []filter{},
		keys:    nil,
		data:    d,
	}
}

func (u *update) buildQueryOrID() bson.D {
	if u.id > 0 {
		return bson.D{{Key: "_id", Value: u.id}}
	}
	query := make(bson.D, 0, len(u.filters))
	for i := range u.filters {
		e := u.filters[i].bson().(bson.E)
		query = append(query, e)
	}
	return query
}

func (u *update) buildObjects() []any {
	result := make([]any, 0, len(u.data))
	for i := range u.data {
		result = append(result, u.valueWrapper(u.data[i], i))
	}
	return result
}

func (u *update) valueWrapper(v any, index int) any {
	if len(u.keys) == 0 || index >= len(u.keys) {
		return v
	}

	return bson.M{u.keys[index]: v}
}

func (u *update) WithKeys(keys ...string) *update {
	u.keys = keys
	return u
}

func (u *update) ID(id uint64) *update {
	if len(u.data) == 1 {
		u.id = id
	}
	return u
}

func (u *update) Exists(key string, exists bool) *update {
	return u.appendFilters("$exists", key, exists)
}

func (u *update) And(key string, value any) *update {
	return u.appendFilters("", key, value)
}

func (u *update) Equal(key string, value any) *update {
	return u.appendFilters("$eq", key, value)
}

func (u *update) NotEqual(key string, value ...any) *update {
	return u.appendFilters("$ne", key, value...).Exists(key, true)
}

func (u *update) Greater(key string, value any) *update {
	return u.appendFilters("$gt", key, value)
}

func (u *update) GreaterOrEqual(key string, value any) *update {
	return u.appendFilters("$gte", key, value)
}

func (u *update) Less(key string, value any) *update {
	return u.appendFilters("$lt", key, value)
}

func (u *update) LessOrEqual(key string, value any) *update {
	return u.appendFilters("$lte", key, value)
}

func (u *update) Bitwise(key string, value any) *update {
	return u.appendFilters("$bitsAllSet", key, value)
}

func (u *update) Contain(key string, value ...any) *update {
	if len(value) == 0 {
		return u
	}
	return u.appendFilters("$all", key, value...)
}

func (u *update) In(key string, value ...any) *update {
	if len(value) == 0 {
		return u
	}
	return u.appendFilters("$in", key, value...)
}

func (u *update) NotIn(key string, value ...any) *update {
	if len(value) == 0 {
		return u
	}
	return u.appendFilters("$nin", key, value...)
}

func (u *update) appendFilters(operation, key string, value ...any) *update {
	u.filters = append(u.filters, newFilter(operation, key, value...))
	return u
}
