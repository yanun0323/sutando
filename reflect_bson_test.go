package sutando

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStructBson(t *testing.T) {
	d := mockData()
	assert.NotNil(t, d)

	bs := bsonEncoder(&d, "Test1", false)
	assert.NotNil(t, bs)
}
