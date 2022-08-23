package example

import (
	"context"

	"github.com/shopspring/decimal"
)

type ExampleInsertStruct struct {
	UserName string          `bson:"user_name"`
	Amount   decimal.Decimal `bson:"amount"`
	Ignore   string          `bson:"-"`
}

func (e *Example) InsertOne(ctx context.Context, obj ExampleInsertStruct) error {
	insert := e.Collection("CollectionName").Insert(&obj)
	_, _, err := e.ExecInsert(ctx, insert)
	if err != nil {
		return err
	}

	insertWithKey := e.Collection("CollectionName").Insert(&obj)
	_, _, err = e.ExecInsert(ctx, insertWithKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *Example) InsertMany(ctx context.Context, objs []ExampleInsertStruct) error {
	insert := e.Collection("CollectionName").Insert(&objs)
	_, _, err := e.ExecInsert(ctx, insert)
	if err != nil {
		return err
	}

	insertWithKey := e.Collection("CollectionName").Insert(&objs)
	_, _, err = e.ExecInsert(ctx, insertWithKey)
	if err != nil {
		return err
	}
	return nil
}
