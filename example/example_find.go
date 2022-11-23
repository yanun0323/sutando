package example

import (
	"context"

	"github.com/shopspring/decimal"
)

type ExampleFindStruct struct {
	UserName string          `bson:"user_name"`
	Amount   decimal.Decimal `bson:"amount"`
	Ignore   string          `bson:"-"`
}

func (e *Example) FindOne(ctx context.Context) (ExampleFindStruct, error) {
	result := ExampleFindStruct{}
	query := e.Collection("CollectionName").Find().
		Equal("user_name", "Yanun").
		Greater("amount", decimal.RequireFromString("300"))
	if err := e.ExecFind(ctx, query, &result); err != nil {
		return ExampleFindStruct{}, err
	}
	return result, nil
}

func (e *Example) FindMany(ctx context.Context) ([]ExampleFindStruct, error) {
	result := []ExampleFindStruct{}
	query := e.Collection("CollectionName").Find().
		Equal("user_name", "Yanun").
		Greater("amount", decimal.RequireFromString("300"))
	if err := e.ExecFind(ctx, query, &result); err != nil {
		return nil, err
	}
	return result, nil
}
