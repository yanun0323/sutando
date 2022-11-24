package sutando

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yanun0323/pkg/logs"
)

type dbSuite struct {
	suite.Suite
	db  DB
	ctx context.Context
	l   *logs.Logger
}

func (su *dbSuite) SetupSuite() {
	ctx := context.Background()
	s, err := NewDB(ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "sutando",
		AdminAuth: true,
		Pem:       "",
	})
	su.Require().Nil(err)
	su.Require().NotNil(s)
	su.db = s
	su.ctx = ctx
	su.l = logs.New("dbSuite", 2)
}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (su dbSuite) Test_CRUD() {
	{
		data := mockData()
		insOne := su.db.Collection("TestCRUD").Insert(&data)
		resultOne, _, err := su.db.ExecInsert(su.ctx, insOne)
		su.Assert().Nil(err)
		su.Assert().NotNil(resultOne)
		su.l.Debug("insert one ID: ", resultOne.InsertedID)

		insMany := su.db.Collection("TestCRUD").Insert(&data, &data, &data)
		_, resultMany, err := su.db.ExecInsert(su.ctx, insMany)
		su.Assert().Nil(err)
		su.Assert().NotNil(resultMany)
		su.l.Debug("insert count: ", len(resultMany.InsertedIDs))

	}

	var One testStruct
	{
		queryOneFist := su.db.Collection("TestCRUD").Find().First()
		su.Nil(su.db.ExecFind(su.ctx, queryOneFist, &One))
		su.NotEmpty(One)

		queryOneFistFailed := su.db.Collection("TestCRUD").Find()
		su.Error(su.db.ExecFind(su.ctx, queryOneFistFailed, &One))
	}

	{
		var a []testStruct

		query := su.db.Collection("TestCRUD").Find()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find all count: ", len(a))
	}

	{
		var a testStruct

		query := su.db.Collection("TestCRUD").Find().Equal("structName", "Yanun").First()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)

		query = su.db.Collection("TestCRUD").Find().Equal("structName", "Yanun")
		err = su.db.ExecFind(su.ctx, query, &a)
		su.Error(err)
	}

	{
		var a []testStruct

		query := su.db.Collection("TestCRUD").Find().Contain("arrTest", 1, 3, 5).First()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find contain first count: ", len(a))

		query = su.db.Collection("TestCRUD").Find().Contain("arrTest", 1, 3, 5)
		err = su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find contain count: ", len(a))
	}

	{
		data := mockData()
		data.StructName = "Vin"
		update := su.db.Collection("TestCRUD").UpdateWith(&data).Equal("structName", "Yanun").First()
		result, err := su.db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("update first count: ", result.ModifiedCount)

		update = su.db.Collection("TestCRUD").UpdateWith(&data).Equal("structName", "Yanun")
		result, err = su.db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("update count: ", result.ModifiedCount)
	}

	{
		query := su.db.Collection("TestCRUD").Delete().Equal("structName", "Vin").First()
		result, err := su.db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("delete count: ", result.DeletedCount)

		query = su.db.Collection("TestCRUD").Delete().Equal("structName", "Vin")
		result, err = su.db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("delete count: ", result.DeletedCount)
	}
}
