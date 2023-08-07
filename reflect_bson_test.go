package sutando

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testGetStructBsonStruct struct {
	FirstField  string
	SecondField string `bson:"secondField"`
	ThirdField  string `bson:"ThirdField"`
	FourthField string
	FifthField  string `bson:"-"`
}

func TestGetStructBson(t *testing.T) {
	{
		d := mockData()
		assert.NotNil(t, d)

		bs := bsonEncoder(&d, "Test1")
		assert.NotNil(t, bs)

	}

	{
		expected := map[string]bool{"firstField": true, "secondField": true, "ThirdField": true, "fourthField": true}
		d := testGetStructBsonStruct{
			FirstField:  "1",
			SecondField: "2",
			ThirdField:  "3",
			FourthField: "4",
			FifthField:  "5",
		}
		hash := bsonEncoder(&d, "")
		assert.NotNil(t, hash)

		for k := range hash {
			assert.True(t, expected[k], fmt.Sprintf("mismatch, %s", k))
		}
	}
}
