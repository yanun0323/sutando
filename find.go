package sutando

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type find struct {
	q querying
}

func newFind(collection *mongo.Collection) finding {
	return &find{
		q: newQuery(collection),
	}
}

func (f *find) Exists(key string, exists bool) finding {
	f.q.Exists(key, exists)
	return f
}

func (f *find) And(key string, value any) finding {
	f.q.And(key, value)
	return f
}

func (f *find) Equal(key string, value any) finding {
	f.q.Equal(key, value)
	return f
}

func (f *find) NotEqual(key string, value ...any) finding {
	f.q.NotEqual(key, value)
	return f
}

func (f *find) Greater(key string, value any) finding {
	f.q.Greater(key, value)
	return f
}

func (f *find) GreaterOrEqual(key string, value any) finding {
	f.q.GreaterOrEqual(key, value)
	return f
}

func (f *find) Less(key string, value any) finding {
	f.q.Less(key, value)
	return f
}

func (f *find) LessOrEqual(key string, value any) finding {
	f.q.LessOrEqual(key, value)
	return f
}

func (f *find) Bitwise(key string, value any) finding {
	f.q.Bitwise(key, value)
	return f
}

func (f *find) Contain(key string, value ...any) finding {
	f.q.Contain(key, value...)
	return f
}

func (f *find) In(key string, value ...any) finding {
	f.q.In(key, value...)
	return f
}

func (f *find) NotIn(key string, value ...any) finding {
	f.q.NotIn(key, value...)
	return f
}

func (f *find) First() finding {
	f.q.First()
	return f
}

func (f *find) Exec(ctx context.Context, p any) error {
	if reflect.TypeOf(p).Kind() != reflect.Pointer {
		return errors.New("object to find should be a pointer")
	}
	kind := reflect.TypeOf(p).Elem().Kind()
	if kind == reflect.Array {
		return errors.New("object to find cannot be an array")
	}

	if kind != reflect.Slice {
		if !f.q.isOne() {
			return errors.New("find too many results! use query with 'First' to find one result")
		}
		return f.execFindOne(ctx, p)
	}

	if !f.q.isOne() {
		return f.execFindMany(ctx, f, p)
	}

	obj := reflect.New(reflect.TypeOf(p).Elem().Elem())
	err := f.execFindOne(ctx, obj.Interface())
	if err != nil {
		return err
	}
	sli := reflect.Append(reflect.ValueOf(p).Elem(), obj.Elem())
	reflect.ValueOf(p).Elem().Set(sli)
	return nil
}

func (f *find) execFindOne(ctx context.Context, p any) error {
	result := f.q.col().FindOne(ctx, f.q.build())
	err := result.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

func (f *find) execFindMany(ctx context.Context, q finding, p any) error {
	cursor, err := f.q.col().Find(ctx, f.q.build())
	if err != nil {
		return err
	}
	err = cursor.All(ctx, p)
	if err != nil {
		return err
	}
	return nil
}
