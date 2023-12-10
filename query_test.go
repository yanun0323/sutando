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
	_, err := su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
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
	su.NoError(su.invokeQuery(col.query().Exists("structName", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Exists("inner.name", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Exists("mapTest.1", true), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Exists("notExistField", false), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Exists("structName", false), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Exists("notExistField", true), &a))
	su.Empty(a)

	// And
	su.NoError(su.invokeQuery(col.query().And("structName", "Yanun"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().And("inner.name", "inner"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().And("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().And("structName", "Peter"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().And("notExistField", "Yanun"), &a))
	su.Empty(a)

	// Equal
	su.NoError(su.invokeQuery(col.query().Equal("structName", "Yanun"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Equal("inner.name", "inner"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Equal("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Equal("structName", "Peter"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Equal("notExistField", "Yanun"), &a))
	su.Empty(a)

	// NotEqual
	su.NoError(su.invokeQuery(col.query().NotEqual("structName", "Yanun"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().NotEqual("inner.name", "inner"), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().NotEqual("mapTest.1", 2), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().NotEqual("structName", "Peter"), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotEqual("notExistField", "Yanun"), &a))
	su.Empty(a)

	// Greater
	su.NoError(su.invokeQuery(col.query().Greater("structAge", 11), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Greater("inner.ohTheAge", 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Greater("mapTest.1", 1), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Greater("structAge", 27), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Greater("inner.ohTheAge", 10), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Greater("structAge", 50), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Greater("inner.ohTheAge", 150), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Greater("mapTest.1", 10), &a))
	su.Empty(a)

	// GreaterOrEqual
	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("structAge", 50), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("inner.ohTheAge", 150), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().GreaterOrEqual("mapTest.1", 10), &a))
	su.Empty(a)

	// Less
	su.NoError(su.invokeQuery(col.query().Less("structAge", 30), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Less("inner.ohTheAge", 18), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Less("mapTest.1", 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Less("structAge", 27), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Less("inner.ohTheAge", 10), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Less("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Less("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Less("mapTest.1", 0), &a))
	su.Empty(a)

	// LessOrEqual
	su.NoError(su.invokeQuery(col.query().LessOrEqual("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().LessOrEqual("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().LessOrEqual("mapTest.1", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().LessOrEqual("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().LessOrEqual("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().LessOrEqual("mapTest.1", 0), &a))
	su.Empty(a)

	// Contain
	su.NoError(su.invokeQuery(col.query().Contain("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Contain("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Contain("arrTest", 2), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().Contain("structAge", 5), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Contain("inner.ohTheAge", 1), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Contain("mapTest.1", 0), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().Contain("arrTest", 6), &a))
	su.Empty(a)

	// In
	su.NoError(su.invokeQuery(col.query().In("structAge", 27), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("inner.ohTheAge", 10), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("inner.ohTheAge", []int{10}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("arrTest", []int{1, 2, 3}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("arrTest", 4, 5), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("arrTest", 1, 3, 6, 7), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().In("arrTest", 6, 7, 8), &a))
	su.Empty(a)

	// NotIn
	su.NoError(su.invokeQuery(col.query().NotIn("structAge", 20), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotIn("inner.ohTheAge", 11), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotIn("inner.ohTheAge", []int{11}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotIn("arrTest", []int{6, 7, 8}), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotIn("arrTest", 0, 6), &a))
	su.Equal(3, len(a))

	su.NoError(su.invokeQuery(col.query().NotIn("arrTest", 5, 6, 7, 8), &a))
	su.Empty(a)

	su.NoError(su.invokeQuery(col.query().NotIn("arrTest", []int{1, 2, 3}), &a))
	su.Empty(a)
}
