package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type insertSuite struct {
	dbSuite
}

func TestInsert(t *testing.T) {
	suite.Run(t, new(insertSuite))
}

func (su *insertSuite) AfterTest(suiteName, testName string) {
	q := su.db.Collection("insert_suite").Delete()
	_, err := su.db.ExecDelete(su.ctx, q)
	su.Require().Nil(err)
}

func (su insertSuite) Test_Insert() {
	data := mockData()
	data.StructName = "Yanun"
	insertOne := su.db.Collection("insert_suite").Insert(&data)
	_, _, err := su.db.ExecInsert(su.ctx, insertOne)
	su.Require().Nil(err)

	insertMany := su.db.Collection("insert_suite").Insert(&data, &data, &data)
	_, _, err = su.db.ExecInsert(su.ctx, insertMany)
	su.Require().Nil(err)
}
