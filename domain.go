package sutando

type query interface {
	ID(id uint64) find
	Exists(key string, exists bool) find
	And(key string, value any) find
	Equal(key string, value any) find
	NotEqual(key string, value ...any) find
	Greater(key string, value any) find
	GreaterOrEqual(key string, value any) find
	Less(key string, value any) find
	LessOrEqual(key string, value any) find
	Bitwise(key string, value any) find
	Contain(key string, value ...any) find
	In(key string, value ...any) find
	NotIn(key string, value ...any) find
}
