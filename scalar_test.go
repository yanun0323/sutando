package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type scalarSuite struct {
	baseSuite

	db  DB
	col string
}

func TestScalar(t *testing.T) {
	suite.Run(t, new(scalarSuite))
}

func (su *scalarSuite) SetupSuite() {
	su.baseSuite.SetupSuite()
	su.col = "find_suite"
}

func (su *scalarSuite) BeforeTest(suiteName, testName string) {
	su.db = su.initDB()
	data := mockData()
	data.StructName = "Yanun"
	_, _, err := su.db.Collection(su.col).Insert(&data, &data, &data).Exec(su.ctx)
	su.Require().NoError(err)
}

func (su *scalarSuite) AfterTest(suiteName, testName string) {
	_, err := su.db.Collection(su.col).Delete().Exec(su.ctx)
	su.Require().NoError(err)
	su.Require().NoError(su.db.Disconnect(su.ctx))
}

func (su *scalarSuite) TestScalarGood() {
	col := su.db.Collection(su.col)

	{
		c, err := col.Scalar().NotIn("structAge", 20).Count(su.ctx)
		su.Require().NoError(err)
		su.EqualValues(3, c)
	}
}
