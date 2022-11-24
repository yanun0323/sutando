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

func (b builder) UpdateWith(p any) update {
	return newUpdate(b.col, p)
}

func (b builder) Update() update {
	return newUpdate(b.col, nil)
}

func (b builder) Find() query {
	return newFind(b.col)
}

func (b builder) Delete() query {
	return newFind(b.col)
}
