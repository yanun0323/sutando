package sutando

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type builderSuite struct {
	dbSuite
}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(builderSuite))
}

func (b *builderSuite) SetupSuite() {
	// TODO: [Yanun] Test
}
