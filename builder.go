package sutando

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

// TODO: custom coder
// withCoder sets custom coder of types/structure.
// func (b builder) withCoder() *builder {
// 	return &b
// }

// Insert provides the operation to insert documents.
//
//	// insert one
//	result, err := db.Collection("col_name").
//		Insert(&obj).
//		Exec(ctx)
//
//	// insert many
//	result, err := db.Collection("col_name").
//		Insert(&obj1, &obj2, &obj3).
//		Exec(ctx)
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, newBsonEncoder(), p...)
}

// Update provides the operation to update documents with structure. (updates all fields, including empty fields)
//
//	// update one
//	result, err := db.Collection("col_name").
//		UpdateWith(&obj).
//		Equal("name", "yanun").
//		Exists("age", true).
//		First().
//		Exec(ctx, true)
//
//	// update many
//	result, err := db.Collection("col_name").
//		UpdateWith(&obj).
//		Equal("name", "yanun").
//		Exists("age", true).
//		Exec(ctx, true)
func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, newBsonEncoder(), p)
}

// Update provides the operation to update documents with set.
//
//	// update one
//	result, err := db.Collection("col_name").
//		Update().
//		Equal("name", "yanun").
//		Exists("age", true).
//		Set("name", "new_name").
//		First().
//		Exec(ctx, true)
//
//	// update many
//	result, err := db.Collection("col_name").
//		Update().
//		Equal("name", "yanun").
//		Exists("age", true).
//		Set("name", "new_name").
//		Exec(ctx, true)
func (b builder) Update() updating {
	return newUpdate(b.col, newBsonEncoder(), nil)
}

// Find provides the operation to find documents.
//
//	// find one
//	err := db.Collection("col_name").
//		Find().
//		Equal("name", "yanun").
//		Exists("age", true).
//		First().
//		Exec(ctx, &elem)
//
//	// find many
//	err := db.Collection("col_name").
//		Find().
//		Equal("name", "yanun").
//		Exists("age", true).
//		Exec(ctx, &elems)
func (b builder) Find() finding {
	return newFind(b.col)
}

// Delete provides the operation to delete documents.
//
//	// delete one
//	err := db.Collection("col_name").
//		Delete().
//		Equal("name", "yanun").
//		Exists("age", true).
//		First().
//		Exec(ctx)
//
//	// delete many
//	err := db.Collection("col_name").
//		Delete().
//		Equal("name", "yanun").
//		Exists("age", true).
//		Exec(ctx)
func (b builder) Delete() deleting {
	return newDelete(b.col)
}

// Scalar provides the operation to scalars documents.
//
//	// count
//	count, err := db.Collection("col_name").
//		Scalar().
//		Equal("name", "yanun").
//		Exists("age", true).
//		Count(ctx, "optional_index1", "optional_index2")
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
