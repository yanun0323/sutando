package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type insertSuite struct {
	baseSuite

	db DB
}

func TestInsert(t *testing.T) {
	suite.Run(t, new(insertSuite))
}

func (su *insertSuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
}

func (su *insertSuite) AfterTest(suiteName, testName string) {
	q := su.db.Collection("insert_suite").Delete()
	_, err := su.db.ExecDelete(su.ctx, q)
	su.Require().NoError(err)
}

func (su *insertSuite) Test_Insert() {
	data := mockData()
	data.StructName = "Yanun"
	insertOne := su.db.Collection("insert_suite").Insert(&data)
	_, _, err := su.db.ExecInsert(su.ctx, insertOne)
	su.Require().NoError(err)

	insertMany := su.db.Collection("insert_suite").Insert(&data, &data, &data)
	_, _, err = su.db.ExecInsert(su.ctx, insertMany)
	su.Require().NoError(err)
}
