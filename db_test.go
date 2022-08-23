package sutando

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type dbSuite struct {
	suite.Suite
	db  DB
	Ctx context.Context
}

func (su *dbSuite) SetupSuite() {
	ctx := context.Background()
	s, err := NewDB(ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "Sutando",
		AdminAuth: true,
		Pem:       "",
	})
	su.Require().Nil(err)
	su.Require().NotNil(s)
	su.db = s
	su.Ctx = ctx
}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (su dbSuite) Test_ExecInsertOne() {
	data := mockData()

	ins := su.db.Collection("TestOne").Insert(&data)
	result, n, err := su.db.ExecInsert(su.Ctx, ins)
	su.Assert().Nil(n)
	su.Assert().Nil(err)
	su.Assert().NotNil(result)

	data.Name = "NotYanun"
	ins = su.db.Collection("TestOne").Insert(&data)
	result, n, err = su.db.ExecInsert(su.Ctx, ins)
	su.Assert().Nil(n)
	su.Assert().Nil(err)
	su.Assert().NotNil(result)
}

func (su dbSuite) Test_ExecInsertMany() {
	data := mockData()

	ins := su.db.Collection("TestMany").Insert(&data, &data)
	n, result, err := su.db.ExecInsert(su.Ctx, ins)

	su.Assert().Nil(n)
	su.Assert().Nil(err)
	su.Assert().NotNil(result)
}

func (su dbSuite) Test_ExecQueryOne_None_Good() {
	query := su.db.Collection("TestOne").Find()
	var a testStruct
	err := su.db.ExecQuery(su.Ctx, query, &a)
	su.Nil(err)
	fmt.Printf("%+v\n", a)
}

func (su dbSuite) Test_ExecQueryOne_Equal_Good() {
	query := su.db.Collection("TestOne").Find().Equal("nameName", "NotYanun")
	var a testStruct
	err := su.db.ExecQuery(su.Ctx, query, &a)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	fmt.Printf("%+v\n", a)
}

func (su dbSuite) Test_ExecQueryMany() {
	query := su.db.Collection("TestMany").Find().Contain("arr", 1, 3, 5)
	var a []testStruct
	err := su.db.ExecQuery(su.Ctx, query, &a)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	fmt.Printf("%+v\n", a)
}

func (su dbSuite) Test_UpdateOne_Equal_Good() {
	data := mockData()
	data.Name = "Vin"
	update := su.db.Collection("TestMany").Update(&data).Equal("nameName", "Yanun")
	result, err := su.db.ExecUpdate(su.Ctx, update, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	fmt.Println("update count: ", result.ModifiedCount)
}
