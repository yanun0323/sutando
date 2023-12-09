package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

// TODO: Implement me
// withCoder sets custom coder of types/structure.
// func (b builder) withCoder() *builder {
// 	return &b
// }

// Insert inserts document.
//
//	err := db.Collection("col_name").Insert(&obj).Exec(ctx)
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, newBsonEncoder(), p...)
}

// Update updates document with model. (updates all fields, including empty fields)
func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, newBsonEncoder(), p)
}

// Update updates document with set.
func (b builder) Update() updating {
	return newUpdate(b.col, newBsonEncoder(), nil)
}

// Find finds document.
func (b builder) Find() finding {
	return newFind(b.col)
}

// Delete deletes document.
func (b builder) Delete() deleting {
	return newDelete(b.col)
}

// Scalar scalars document.
func (b builder) Scalar() querying {
	return newQuery(b.col)
}

// Drop drops this collection from database.
func (b builder) Drop(ctx context.Context) error {
	return b.col.Drop(ctx)
}

// query is for package testing.
func (b builder) query() querying {
	return newQuery(b.col)
}
