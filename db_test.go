package sutando

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type baseSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *baseSuite) SetupSuite() {
	ctx := context.Background()
	su.ctx = ctx
}

func (su *baseSuite) initDB() DB {
	s, err := NewDB(su.ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "sutando",
		AdminAuth: true,
		Pem:       "",
		Srv:       false,
	})
	su.Require().NoError(err)
	su.Require().NotNil(s)
	return s
}

func (su *baseSuite) checkEmpty(col string) {
	db := su.initDB()
	var result []map[string]interface{}
	err := db.Collection(col).Find().Exec(su.ctx, &result)
	su.Require().NoError(err)
	su.Require().Empty(result)
}

type dbSuite struct {
	baseSuite

	col string
}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (su *dbSuite) SetupSuite() {
	su.baseSuite.SetupSuite()
	su.col = "db_suite"
}

func (su *dbSuite) TestSrvGood() {
	if os.Getenv("TEST_SRV") != "1" {
		return
	}
	s, err := NewDB(su.ctx, Conn{
		Username:  "test",
		Password:  "test",
		Host:      "sutando.mongodb.net",
		DB:        "test",
		AdminAuth: true,
		Srv:       true,
		OptionHandler: func(client *options.ClientOptions) {
			client.SetConnectTimeout(3 * time.Second)
			client.SetTimeout(3 * time.Second)
		},
	})
	su.Require().NoError(err)
	su.Require().NotNil(s)

	var result map[string]interface{}
	su.Require().NoError(s.Collection(su.col).Find().First().Exec(su.ctx, &result))
	su.T().Log(len(result))
}

func (su *dbSuite) TestNewDBFromGood() {
	client, err := mongo.Connect(su.ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s?authenticationDatabase=admin",
		"test", "test", "localhost:27017", "sutando")))
	su.NoError(err)
	su.NotNil(client)
	db := NewDBFromMongo(su.ctx, client, "sutando")
	su.NotNil(db)
}

func (su *dbSuite) TestDisconnectGood() {
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

	su.NoError(db.Disconnect(su.ctx))
}

func (su *dbSuite) TestCRUDGood() {
	db := su.initDB()
	{
		data := mockData()
		resultOne, _, err := db.Collection(su.col).Insert(&data).Exec(su.ctx)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultOne)
		su.T().Log("insert one ID: ", resultOne.InsertedID)

		_, resultMany, err := db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultMany)
		su.T().Log("insert count: ", len(resultMany.InsertedIDs))

	}

	{
		var a testStruct

		su.Nil(db.Collection(su.col).Find().Equal("structName", "Yanun").First().Exec(su.ctx, &a))
		su.NotEmpty(a)

		su.Error(db.Collection(su.col).Find().Exec(su.ctx, &a))
	}

	{
		var a []testStruct

		err := db.Collection(su.col).Find().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find all count: ", len(a))
	}

	{
		var a testStruct

		err := db.Collection(su.col).Find().Equal("structName", "Yanun").First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)

		err = db.Collection(su.col).Find().Equal("structName", "Yanun").Exec(su.ctx, &a)
		su.Error(err)
	}

	{
		var a []testStruct

		err := db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5).First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find contain first count: ", len(a))

		err = db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5).Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find contain count: ", len(a))
	}

	{
		data := mockData()
		data.StructName = "Vin"
		result, err := db.Collection(su.col).UpdateWith(&data).Equal("structName", "Yanun").First().Exec(su.ctx, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("update first count: ", result.ModifiedCount)

		result, err = db.Collection(su.col).UpdateWith(&data).Equal("structName", "Yanun").Exec(su.ctx, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("update count: ", result.ModifiedCount)
	}

	{
		result, err := db.Collection(su.col).Delete().Equal("structName", "Vin").First().Exec(su.ctx)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("delete count: ", result.DeletedCount)

		result, err = db.Collection(su.col).Delete().Equal("structName", "Vin").Exec(su.ctx)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("delete count: ", result.DeletedCount)
	}

	{
		resultOne, _, err := db.Collection(su.col).Insert(_defaultSettings).Exec(su.ctx)
		su.NoError(err)
		su.NotEmpty(resultOne.InsertedID)
	}

	{
		var a Setting
		err := db.Collection(su.col).Find().Equal("exchangeSettings.OkCoin.buyEnable", true).First().Exec(su.ctx, &a)
		su.NoError(err)
		su.NotZero(a)
	}

	{
		result, err := db.Collection(su.col).Delete().Equal("exchangeSettings.OkCoin.sellEnable", true).Exec(su.ctx)
		su.NoError(err)
		su.NotZero(result.DeletedCount)
	}

	su.checkEmpty(su.col)
	su.NoError(db.Disconnect(su.ctx))
}

func (su *dbSuite) TestOldGood() {
	db := su.initDB()
	{
		data := mockData()
		insOne := db.Collection(su.col).Insert(&data)
		resultOne, _, err := db.ExecInsert(su.ctx, insOne)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultOne)
		su.T().Log("insert one ID: ", resultOne.InsertedID)

		insMany := db.Collection(su.col).Insert(&data, &data, &data)
		_, resultMany, err := db.ExecInsert(su.ctx, insMany)
		su.Assert().NoError(err)
		su.Assert().NotNil(resultMany)
		su.T().Log("insert count: ", len(resultMany.InsertedIDs))

	}

	{
		var a testStruct

		queryOneFist := db.Collection(su.col).Find().First()
		su.Nil(db.ExecFind(su.ctx, queryOneFist, &a))
		su.NotEmpty(a)

		queryOneFistFailed := db.Collection(su.col).Find()
		su.Error(db.ExecFind(su.ctx, queryOneFistFailed, &a))
	}

	{
		var a []testStruct

		query := db.Collection(su.col).Find()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find all count: ", len(a))
	}

	{
		var a testStruct

		query := db.Collection(su.col).Find().Equal("structName", "Yanun").First()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)

		query = db.Collection(su.col).Find().Equal("structName", "Yanun")
		err = db.ExecFind(su.ctx, query, &a)
		su.Error(err)
	}

	{
		var a []testStruct

		query := db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5).First()
		err := db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find contain first count: ", len(a))

		query = db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5)
		err = db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.T().Log("find contain count: ", len(a))
	}

	{
		data := mockData()
		data.StructName = "Vin"
		update := db.Collection(su.col).UpdateWith(&data).Equal("structName", "Yanun").First()
		result, err := db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("update first count: ", result.ModifiedCount)

		update = db.Collection(su.col).UpdateWith(&data).Equal("structName", "Yanun")
		result, err = db.ExecUpdate(su.ctx, update, false)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("update count: ", result.ModifiedCount)
	}

	{
		query := db.Collection(su.col).Delete().Equal("structName", "Vin").First()
		result, err := db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("delete count: ", result.DeletedCount)

		query = db.Collection(su.col).Delete().Equal("structName", "Vin")
		result, err = db.ExecDelete(su.ctx, query)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.T().Log("delete count: ", result.DeletedCount)
	}

	su.NoError(db.Disconnect(su.ctx))
}
