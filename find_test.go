package sutando

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type findSuite struct {
	baseSuite

	db  DB
	col string
}

func TestFind(t *testing.T) {
	suite.Run(t, new(findSuite))
}

func (su *findSuite) SetupSuite() {
	su.baseSuite.SetupSuite()
	su.col = "find_suite"
}

func (su *findSuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
	data := mockData()
	data.StructName = "Yanun"
	_, _, err := su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *findSuite) AfterTest(suiteName, testName string) {
	_, err := su.db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NoError(su.db.Disconnect(su.ctx))
}

func (su *findSuite) Test_Find_Good() {
	{
		var a testStruct
		su.NoError(su.db.Collection(su.col).Find().First().Exec(su.ctx, &a))
		su.NotEmpty(a)

		su.NoError(su.db.Collection(su.col).Find().First().Exec(su.ctx, &a))
	}

	{
		var a []testStruct

		err := su.db.Collection(su.col).Find().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))

		su.NoError(su.db.Collection(su.col).Find().First().Exec(su.ctx, &a))
	}

	{
		var a testStruct

		err := su.db.Collection(su.col).Find().Equal("structName", "Yanun").First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
	}

	{
		var a []testStruct

		err := su.db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5).First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(1, len(a))

		err = su.db.Collection(su.col).Find().Contain("arrTest", 1, 3, 5).Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))
	}
}

func (su *findSuite) Test_Find_Bad() {
	{
		var a testStruct
		su.Error(su.db.Collection(su.col).Find().Exec(su.ctx, &a))
	}

	{
		var a testStruct
		su.Error(su.db.Collection(su.col).Find().Equal("structName", "Yanun").Exec(su.ctx, &a))
	}
}
