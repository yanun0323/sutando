package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

/*
SetCustomEncodeTypes sets custom encode types/structure
*/
func (b builder) SetCustomEncodeTypes() *builder {
	return &b
}

/*
Insert inserts data
*/
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, newBsonEncoder(), p...)
}

/*
Update updates data with model (Will update all fields including empty fields)
*/
func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, newBsonEncoder(), p)
}

/*
Update updates data with set
*/
func (b builder) Update() updating {
	return newUpdate(b.col, newBsonEncoder(), nil)
}

/*
Find finds data
*/
func (b builder) Find() finding {
	return newFind(b.col)
}

/*
Delete deletes data
*/
func (b builder) Delete() deleting {
	return newDelete(b.col)
}

/*
Scalar scalars data
*/
func (b builder) Scalar() scalaring {
	return newScalar(b.col)
}

/*
Drop drops this collection from database
*/
func (b builder) Drop(ctx context.Context) error {
	return b.col.Drop(ctx)
}

/*
query is for package testing
*/
func (b builder) query() querying {
	return newQuery(b.col)
}
