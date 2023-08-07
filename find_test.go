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

func (su *findSuite) TestFindGood() {
	col := su.db.Collection(su.col)

	{
		var a testStruct
		su.NoError(col.Find().First().Exec(su.ctx, &a))
		su.NotEmpty(a)

		su.NoError(col.Find().First().Exec(su.ctx, &a))
	}

	{
		var a []testStruct

		err := col.Find().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))

		su.NoError(col.Find().First().Exec(su.ctx, &a))
	}

	{
		var a testStruct

		err := col.Find().Equal("structName", "Yanun").First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
	}

	{
		var a []testStruct

		err := col.Find().Contain("arrTest", 1, 3, 5).First().Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(1, len(a))

		err = col.Find().Contain("arrTest", 1, 3, 5).Exec(su.ctx, &a)
		su.True(err == nil || errors.Is(err, ErrNoDocument), err)
		su.NotEmpty(a)
		su.Equal(3, len(a))
	}

	{
		var a testStruct
		su.Error(col.Find().Exec(su.ctx, &a))

		su.Error(col.Find().Equal("structName", "Yanun").Exec(su.ctx, &a))
	}

	{
		var a []testStruct

		// Exists
		su.NoError(col.Find().Exists("structName", true).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Exists("inner.name", true).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Exists("mapTest.1", true).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Exists("notExistField", false).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Exists("structName", false).Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().Exists("notExistField", true).Exec(su.ctx, &a))
		su.Empty(a)

		// And
		su.NoError(col.Find().And("structName", "Yanun").Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().And("inner.name", "inner").Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().And("mapTest.1", 2).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().And("structName", "Peter").Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().And("notExistField", "Yanun").Exec(su.ctx, &a))
		su.Empty(a)

		// Equal
		su.NoError(col.Find().Equal("structName", "Yanun").Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Equal("inner.name", "inner").Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Equal("mapTest.1", 2).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Equal("structName", "Peter").Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().Equal("notExistField", "Yanun").Exec(su.ctx, &a))
		su.Empty(a)

		// NotEqual
		su.NoError(col.Find().NotEqual("structName", "Yanun").Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().NotEqual("inner.name", "inner").Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().NotEqual("mapTest.1", 2).Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().NotEqual("structName", "Peter").Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().NotEqual("notExistField", "Yanun").Exec(su.ctx, &a))
		su.Empty(a)

		// Greater
		su.NoError(col.Find().Greater("structAge", 11).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Greater("inner.ohTheAge", 5).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Greater("mapTest.1", 1).Exec(su.ctx, &a))
		su.Equal(3, len(a))

		su.NoError(col.Find().Greater("structAge", 50).Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().Greater("inner.ohTheAge", 150).Exec(su.ctx, &a))
		su.Empty(a)

		su.NoError(col.Find().Greater("mapTest.1", 10).Exec(su.ctx, &a))
		su.Empty(a)

		// GreaterOrEqual
		// Less
		// LessOrEqual
		// Bitwise
		// Contain
		// In
		// NotIn

	}
}
