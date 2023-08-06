package sutando

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

/*
 */
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, p...)
}

func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, p)
}

func (b builder) Update() updating {
	return newUpdate(b.col, nil)
}

func (b builder) Find() finding {
	return newFind(b.col)
}

func (b builder) Delete() deleting {
	return newDelete(b.col)
}
