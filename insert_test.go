package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type insertSuite struct {
	baseSuite

	db  DB
	col string
}

func TestInsert(t *testing.T) {
	suite.Run(t, new(insertSuite))
}

func (su *insertSuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
	su.col = "insert_suite"
	_, err := su.db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *insertSuite) AfterTest(suiteName, testName string) {
	su.Require().NoError(su.db.Disconnect(su.ctx))
}

func (su *insertSuite) TestInsertGood() {
	data := mockData()
	data.StructName = "Yanun"

	result, err := su.db.Collection(su.col).Insert(&data).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NotNil(result.InsertedIDs)
	su.T().Log(result.InsertedIDs[0])

	result, err = su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NotEmpty(result.InsertedIDs)
	su.T().Log(result.InsertedIDs)

	dataSlice := []testStruct{data, data, data}
	result, err = su.db.Collection(su.col).Insert(dataSlice).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NotEmpty(result.InsertedIDs)
	su.T().Log(result.InsertedIDs)

	{
		result, err := su.db.Collection(su.col).Insert(_defaultSettings).Exec(su.ctx)
		su.Require().NoError(err)
		su.Require().NotNil(result.InsertedIDs[0])
		su.T().Log(result.InsertedIDs[0])

		result, err = su.db.Collection(su.col).Insert(_defaultSettings, _defaultSettings, _defaultSettings).Exec(su.ctx)
		su.Require().NoError(err)
		su.Require().NotEmpty(result.InsertedIDs)
		su.T().Log(result.InsertedIDs)
	}
}
