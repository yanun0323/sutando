package sutando

type IQuery interface {
	ID(id uint64) query
	Exists(key string, exists bool) query
	And(key string, value any) query
	Equal(key string, value any) query
	NotEqual(key string, value ...any) query
	Greater(key string, value any) query
	GreaterOrEqual(key string, value any) query
	Less(key string, value any) query
	LessOrEqual(key string, value any) query
	Bitwise(key string, value any) query
	Contain(key string, value ...any) query
	In(key string, value ...any) query
	NotIn(key string, value ...any) query
}
