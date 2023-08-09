package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type deleteSuite struct {
	baseSuite

	col string
}

func TestDelete(t *testing.T) {
	suite.Run(t, new(deleteSuite))
}

func (su *deleteSuite) SetupSuite() {
	su.baseSuite.SetupSuite()
	su.col = "delete_suite"
}

func (su *deleteSuite) SetupTest() {
	db := su.initDB()
	_, err := db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NoError(db.Disconnect(su.ctx))
}

func (su *deleteSuite) TestDeleteGood() {
	type Data struct {
		Name  string
		Age   int
		Hired bool
	}

	var data = []Data{
		{
			Name:  "Peter",
			Age:   30,
			Hired: true,
		},
		{
			Name:  "Yanun",
			Age:   12,
			Hired: false,
		},
		{
			Name:  "SpiderMan",
			Age:   -1,
			Hired: true,
		},
	}

	db := su.initDB()
	su.insertData(db, data)
	result, err := db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().Equal(int64(3), result.DeletedCount)

	su.insertData(db, data)
	result, err = db.Collection(su.col).Delete().Equal("name", "Yanun").Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().Equal(int64(1), result.DeletedCount)

	result, err = db.Collection(su.col).Delete().Less("age", 0).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().Equal(int64(1), result.DeletedCount)

	result, err = db.Collection(su.col).Delete().Exists("exist", true).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().Equal(int64(0), result.DeletedCount)

	result, err = db.Collection(su.col).Delete().Equal("hired", true).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().Equal(int64(1), result.DeletedCount)

}

func (su *deleteSuite) insertData(db DB, data ...any) {
	{
		_, _, err := db.Collection(su.col).Insert(data...).Exec(su.ctx)
		su.Require().NoError(err)
	}
}
