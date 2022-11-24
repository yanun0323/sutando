package sutando

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type findSuite struct {
	dbSuite
}

func TestFind(t *testing.T) {
	suite.Run(t, new(findSuite))
}

func (su *findSuite) BeforeTest(suiteName, testName string) {
	data := mockData()
	data.StructName = "Yanun"
	q := su.db.Collection("find_suite").Insert(&data, &data, &data)
	_, _, err := su.db.ExecInsert(su.ctx, q)
	su.Require().Nil(err)
}

func (su *findSuite) AfterTest(suiteName, testName string) {
	q := su.db.Collection("find_suite").Delete()
	_, err := su.db.ExecDelete(su.ctx, q)
	su.Require().Nil(err)
}

func (su findSuite) Test_Find() {
	{
		var One testStruct
		queryOneFist := su.db.Collection("find_suite").Find().First()
		su.Nil(su.db.ExecFind(su.ctx, queryOneFist, &One))
		su.NotEmpty(One)

		queryOneFistFailed := su.db.Collection("find_suite").Find()
		su.Error(su.db.ExecFind(su.ctx, queryOneFistFailed, &One))
	}

	{
		var a []testStruct

		query := su.db.Collection("find_suite").Find()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))
	}

	{
		var a testStruct

		query := su.db.Collection("find_suite").Find().Equal("structName", "Yanun").First()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)

		query = su.db.Collection("find_suite").Find().Equal("structName", "Yanun")
		err = su.db.ExecFind(su.ctx, query, &a)
		su.Error(err)
	}

	{
		var a []testStruct

		query := su.db.Collection("find_suite").Find().Contain("arrTest", 1, 3, 5).First()
		err := su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(1, len(a))

		query = su.db.Collection("find_suite").Find().Contain("arrTest", 1, 3, 5)
		err = su.db.ExecFind(su.ctx, query, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))
	}
}
