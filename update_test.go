package sutando

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type updateSuite struct {
	baseSuite

	db DB
}

func TestUpdate(t *testing.T) {
	suite.Run(t, new(updateSuite))
}

func (su *updateSuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
	_, err := su.db.Collection("update_suite").Delete().Exec(su.ctx)
	su.Require().NoError(err)
	data := mockData()
	data.StructName = "Yanun"
	q := su.db.Collection("update_suite").Insert(&data, &data, &data)
	_, _, err = q.Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *updateSuite) AfterTest(suiteName, testName string) {
	su.Require().NoError(su.db.Disconnect(su.ctx))
}

func (su *updateSuite) TestFindGood() {
	data := mockData()
	data.StructName = "Vin"
	data.StructAge = 50
	result, err := su.db.Collection("update_suite").UpdateWith(&data).Equal("structName", "Yanun").First().Exec(su.ctx, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(1), result.ModifiedCount)

	result, err = su.db.Collection("update_suite").UpdateWith(&data).Equal("structName", "Yanun").Exec(su.ctx, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(2), result.ModifiedCount)

	result, err = su.db.Collection("update_suite").Update().Equal("structName", "Vin").Set("structName", "Yanun").First().Exec(su.ctx, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(1), result.ModifiedCount)

	result, err = su.db.Collection("update_suite").Update().Equal("structName", "Vin").Set("structName", "Yanun").Exec(su.ctx, false)
	su.True(err == nil || errors.Is(err, ErrNoDocument), err)
	su.Equal(int64(2), result.ModifiedCount)

	d := testStruct{}
	su.Require().Nil(su.db.Collection("update_suite").Find().Equal("structAge", 50).First().Exec(su.ctx, &d))
	su.Equal(50, d.StructAge)
}
