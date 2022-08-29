package example

import (
	"context"

	"github.com/yanun0323/sutando"
)

type Example struct {
	sutando.DB
}

func NewExample(ctx context.Context) *Example {
	db, err := sutando.NewDB(ctx, sutando.Conn{
		Username:  "example",
		Password:  "example",
		Host:      "localhost",
		Port:      27017,
		DB:        "example",
		AdminAuth: true,
		Pem:       "",
	})
	if err != nil {
		return nil
	}
	return &Example{
		DB: db,
	}
}

func NewExample2(ctx context.Context) *Example {
	db, err := sutando.NewDB(ctx, sutando.Conn{
		Username:  "example",
		Password:  "example",
		Host:      "localhost:27017",
		DB:        "example",
		AdminAuth: true,
		Pem:       "",
	})
	if err != nil {
		return nil
	}
	return &Example{
		DB: db,
	}
}
