package sutando

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type updateSuite struct {
	dbSuite
}

func TestUpdate(t *testing.T) {
	suite.Run(t, new(updateSuite))
}

func (su *updateSuite) BeforeTest(suiteName, testName string) {
	data := mockData()
	data.StructName = "Yanun"
	q := su.db.Collection("update_suite").Insert(&data, &data, &data)
	_, _, err := su.db.ExecInsert(su.ctx, q)
	su.Require().Nil(err)
}

func (su *updateSuite) AfterTest(suiteName, testName string) {
	q := su.db.Collection("update_suite").Delete()
	_, err := su.db.ExecDelete(su.ctx, q)
	su.Require().Nil(err)
}

func (su updateSuite) Test_Find() {
	data := mockData()
	data.StructName = "Vin"
	data.StructAge = 50
	updateOneWithModel := su.db.Collection("update_suite").UpdateWith(&data).Equal("structName", "Yanun").First()
	result, err := su.db.ExecUpdate(su.ctx, updateOneWithModel, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(1), result.ModifiedCount)

	updateManyWithModel := su.db.Collection("update_suite").UpdateWith(&data).Equal("structName", "Yanun")
	result, err = su.db.ExecUpdate(su.ctx, updateManyWithModel, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(2), result.ModifiedCount)

	updateOne := su.db.Collection("update_suite").Update().Equal("structName", "Vin").Set("structName", "Yanun").First()
	result, err = su.db.ExecUpdate(su.ctx, updateOne, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(1), result.ModifiedCount)

	updateMany := su.db.Collection("update_suite").Update().Equal("structName", "Vin").Set("structName", "Yanun")
	result, err = su.db.ExecUpdate(su.ctx, updateMany, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(2), result.ModifiedCount)

	d := testStruct{}
	q := su.db.Collection("update_suite").Find().Equal("structAge", 50).First()
	su.Require().Nil(su.db.ExecFind(su.ctx, q, &d))
	su.Equal(50, d.StructAge)
}
