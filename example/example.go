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
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "Sutando",
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
