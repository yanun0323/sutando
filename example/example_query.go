package example

import (
	"context"

	"github.com/shopspring/decimal"
)

type ExampleQueryStruct struct {
	UserName string          `bson:"user_name"`
	Amount   decimal.Decimal `bson:"amount"`
	Ignore   string          `bson:"-"`
}

func (e *Example) QueryOne(ctx context.Context) (ExampleQueryStruct, error) {
	result := ExampleQueryStruct{}
	query := e.Collection("CollectionName").Find().
		Equal("user_name", "Yanun").
		Greater("amount", decimal.RequireFromString("300"))
	if err := e.ExecQuery(ctx, query, &result); err != nil {
		return ExampleQueryStruct{}, err
	}
	return result, nil
}

func (e *Example) QueryMany(ctx context.Context) ([]ExampleQueryStruct, error) {
	result := []ExampleQueryStruct{}
	query := e.Collection("CollectionName").Find().
		Equal("user_name", "Yanun").
		Greater("amount", decimal.RequireFromString("300"))
	if err := e.ExecQuery(ctx, query, &result); err != nil {
		return nil, err
	}
	return result, nil
}
