package sutando

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

func (b builder) Insert(p ...any) *insert {
	return newInsert(b.col, p...)
}

func (b builder) Update(p ...any) update {
	return newUpdate(b.col, p...)
}

func (b builder) Find() query {
	return newFind(b.col)
}

func (b builder) Delete(p ...any) *delete {
	return newDelete(b.col, p...)
}
