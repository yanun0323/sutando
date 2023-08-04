package sutando

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yanun0323/pkg/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type baseSuite struct {
	suite.Suite
	ctx context.Context
	l   *logs.Logger
}

func (su *baseSuite) SetupSuite() {
	ctx := context.Background()
	su.ctx = ctx
	su.l = logs.New("dbSuite", 2)
}

func (su *baseSuite) initDB() DB {
	s, err := NewDB(su.ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      0,
		DB:        "sutando",
		AdminAuth: true,
		Pem:       "",
		Srv:       false,
	})
	su.Require().NoError(err)
	su.Require().NotNil(s)
	return s
}

type dbSuite struct {
	baseSuite
}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (su *dbSuite) Test_Srv() {
	s, err := NewDB(su.ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "sutando.com",
		Port:      0,
		DB:        "sutando",
		AdminAuth: false,
		Pem:       "",
		Srv:       true,
		OptionHandler: func(opt *options.ClientOptions) {
			opt.SetConnectTimeout(time.Second)
			opt.SetTimeout(time.Second)
		},
	})
	su.Require().NoError(err)
	su.Require().NotNil(s)
}

func (su *dbSuite) Test_NewDBFrom() {
	client, err := mongo.Connect(su.ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s?authenticationDatabase=admin",
		"test", "test", "localhost:27017", "sutando")))
	su.NoError(err)
	su.NotNil(client)
	db := NewDBFromMongo(su.ctx, client, "sutando")
	su.NotNil(db)
}

func (su *dbSuite) Test_Disconnect() {
	db, err := NewDB(su.ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "sutando",
		AdminAuth: true,
		Pem:       "",
	})
	su.Require().NoError(err)
	su.Require().NotNil(db)

	su.Nil(db.Disconnect(su.ctx))
}

func (su *dbSuite) Test_CRUD() {
	db := su.initDB()
	col := "db_suite"
	{
		data := mockData()
		insOne := db.Collection(col).Insert(&data)
		resultOne, _, err := db.ExecInsert(su.ctx, insOne)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultOne)
		su.l.Debug("insert one ID: ", resultOne.InsertedID)

		insMany := db.Collection(col).Insert(&data, &data, &data)
		_, resultMany, err := db.ExecInsert(su.ctx, insMany)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultMany)
		su.l.Debug("insert count: ", len(resultMany.InsertedIDs))

	}

	var One testStruct
	{
		queryOneFist := db.Collection(col).Find().First()
		su.Nil(db.ExecFind(su.ctx, queryOneFist, &One))
		su.NotEmpty(One)

		queryOneFistFailed := db.Collection(col).Find()
		su.Error(db.ExecFind(su.ctx, queryOneFistFailed, &One))
	}

	{
		var a []testStruct

		query := db.Collection(col).Find()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find all count: ", len(a))
	}

	{
		var a testStruct

		query := db.Collection(col).Find().Equal("structName", "Yanun").First()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)

		query = db.Collection(col).Find().Equal("structName", "Yanun")
		err = db.ExecFind(su.ctx, query, &a)
		su.Error(err)
	}

	{
		var a []testStruct

		query := db.Collection(col).Find().Contain("arrTest", 1, 3, 5).First()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find contain first count: ", len(a))

		query = db.Collection(col).Find().Contain("arrTest", 1, 3, 5)
		err = db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.l.Debug("find contain count: ", len(a))
	}

	{
		data := mockData()
		data.StructName = "Vin"
		update := db.Collection(col).UpdateWith(&data).Equal("structName", "Yanun").First()
		result, err := db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("update first count: ", result.ModifiedCount)

		update = db.Collection(col).UpdateWith(&data).Equal("structName", "Yanun")
		result, err = db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("update count: ", result.ModifiedCount)
	}

	{
		query := db.Collection(col).Delete().Equal("structName", "Vin").First()
		result, err := db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("delete count: ", result.DeletedCount)

		query = db.Collection(col).Delete().Equal("structName", "Vin")
		result, err = db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.l.Debug("delete count: ", result.DeletedCount)
	}
}
