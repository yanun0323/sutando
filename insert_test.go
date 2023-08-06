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
}

func (su *insertSuite) AfterTest(suiteName, testName string) {
	_, err := su.db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *insertSuite) Test_Insert_Good() {
	data := mockData()
	data.StructName = "Yanun"

	resultOne, _, err := su.db.Collection(su.col).Insert(&data).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NotNil(resultOne.InsertedID)

	_, resultMany, err := su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NotEmpty(resultMany.InsertedIDs)
}
