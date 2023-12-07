package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scalar struct {
	q querying
}

func newScalar(collection *mongo.Collection) scalaring {
	return &scalar{
		q: newQuery(collection),
	}
}

func (s *scalar) Exists(key string, exists bool) scalaring {
	s.q.Exists(key, exists)
	return s
}

func (s *scalar) And(key string, value any) scalaring {
	s.q.And(key, value)
	return s
}

func (s *scalar) Equal(key string, value any) scalaring {
	s.q.Equal(key, value)
	return s
}

func (s *scalar) NotEqual(key string, value ...any) scalaring {
	s.q.NotEqual(key, value...)
	return s
}

func (s *scalar) Greater(key string, value any) scalaring {
	s.q.Greater(key, value)
	return s
}

func (s *scalar) GreaterOrEqual(key string, value any) scalaring {
	s.q.GreaterOrEqual(key, value)
	return s
}

func (s *scalar) Less(key string, value any) scalaring {
	s.q.Less(key, value)
	return s
}

func (s *scalar) LessOrEqual(key string, value any) scalaring {
	s.q.LessOrEqual(key, value)
	return s
}

func (s *scalar) Bitwise(key string, value any) scalaring {
	s.q.Bitwise(key, value)
	return s
}

func (s *scalar) Contain(key string, value ...any) scalaring {
	s.q.Contain(key, value...)
	return s
}

func (s *scalar) In(key string, value ...any) scalaring {
	s.q.In(key, value...)
	return s
}

func (s *scalar) NotIn(key string, value ...any) scalaring {
	s.q.NotIn(key, value...)
	return s
}

func (s *scalar) Regex(key string, regex string) scalaring {
	s.q.Regex(key, regex)
	return s
}

func (s *scalar) Count(ctx context.Context, index ...string) (int64, error) {
	opt := options.Count()
	if len(index) != 0 && len(index[0]) != 0 {
		opt.SetHint(index[0])
	}
	i, err := s.q.col().CountDocuments(ctx, s.q.build(), opt)
	if err != nil {
		return s.q.col().CountDocuments(ctx, s.q.build(), options.Count())
	}
	return i, nil
}
