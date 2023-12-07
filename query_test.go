package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type querySuite struct {
	baseSuite

	db  DB
	col string
}

func TestQuery(t *testing.T) {
	suite.Run(t, new(querySuite))
}

func (su *querySuite) SetupSuite() {
	su.baseSuite.SetupSuite()
	su.col = "query_suite"
}

func (su *querySuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
	data := mockData()
	data.StructName = "Yanun"
	_, _, err := su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *querySuite) AfterTest(suiteName, testName string) {
	_, err := su.db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NoError(su.db.Disconnect(su.ctx))
}

func (su *querySuite) invokeQuery(q querying, p any) error {
	find := &find{
		q: q,
	}
	return find.Exec(su.ctx, p)
}

func (su *querySuite) TestQueryGood() {
	col := su.db.Collection(su.col)

	var a []testStruct

	// Exists
	su.NoError(su.invokeQuery(col.Query().Exists("structName", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Exists("inner.name", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Exists("mapTest.1", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Exists("notExistField", false), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Exists("structName", false), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Exists("notExistField", true), &a))
	su.Empty(a)

	// And
	su.NoError(su.invokeQuery(col.Query().And("structName", "Yanun"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().And("inner.name", "inner"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().And("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().And("structName", "Peter"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().And("notExistField", "Yanun"), &a))
	su.Empty(a)

	// Equal
	su.NoError(su.invokeQuery(col.Query().Equal("structName", "Yanun"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Equal("inner.name", "inner"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Equal("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Equal("structName", "Peter"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Equal("notExistField", "Yanun"), &a))
	su.Empty(a)

	// NotEqual
	su.NoError(su.invokeQuery(col.Query().NotEqual("structName", "Yanun"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().NotEqual("inner.name", "inner"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().NotEqual("mapTest.1", 2), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().NotEqual("structName", "Peter"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotEqual("notExistField", "Yanun"), &a))
	su.Empty(a)

	// Greater
	su.NoError(su.invokeQuery(col.Query().Greater("structAge", 11), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Greater("inner.ohTheAge", 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Greater("mapTest.1", 1), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Greater("structAge", 27), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Greater("inner.ohTheAge", 10), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Greater("structAge", 50), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Greater("inner.ohTheAge", 150), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Greater("mapTest.1", 10), &a))
	su.Empty(a)

	// GreaterOrEqual
	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("structAge", 50), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("inner.ohTheAge", 150), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().GreaterOrEqual("mapTest.1", 10), &a))
	su.Empty(a)

	// Less
	su.NoError(su.invokeQuery(col.Query().Less("structAge", 30), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Less("inner.ohTheAge", 18), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Less("mapTest.1", 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Less("structAge", 27), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Less("inner.ohTheAge", 10), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Less("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Less("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Less("mapTest.1", 0), &a))
	su.Empty(a)

	// LessOrEqual
	su.NoError(su.invokeQuery(col.Query().LessOrEqual("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().LessOrEqual("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().LessOrEqual("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().LessOrEqual("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().LessOrEqual("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().LessOrEqual("mapTest.1", 0), &a))
	su.Empty(a)

	// Bitwise
	// TODO: Implement me

	// Contain
	su.NoError(su.invokeQuery(col.Query().Contain("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Contain("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Contain("arrTest", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().Contain("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Contain("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Contain("mapTest.1", 0), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().Contain("arrTest", 6), &a))
	su.Empty(a)

	// In
	su.NoError(su.invokeQuery(col.Query().In("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("inner.ohTheAge", []int{10}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("arrTest", []int{1, 2, 3}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("arrTest", 4, 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("arrTest", 1, 3, 6, 7), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().In("arrTest", 6, 7, 8), &a))
	su.Empty(a)

	// NotIn
	su.NoError(su.invokeQuery(col.Query().NotIn("structAge", 20), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotIn("inner.ohTheAge", 11), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotIn("inner.ohTheAge", []int{11}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotIn("arrTest", []int{6, 7, 8}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotIn("arrTest", 0, 6), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.Query().NotIn("arrTest", 5, 6, 7, 8), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.Query().NotIn("arrTest", []int{1, 2, 3}), &a))
	su.Empty(a)

	// Count
	c, err := col.Query().NotIn("structAge", 20).Count(su.ctx)
	su.Require().NoError(err)
	su.EqualValues(3, c)
}
