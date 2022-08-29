package sutando

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTagStruct struct {
	FirstField  string
	SecondField string `bson:"secondField"`
	thirdField  string `bson:"ThirdField"`
	fourthField string
	fifthField  string `bson:"-"`
}

func Test_getTag(t *testing.T) {
	data := &testTagStruct{}
	elem := reflect.ValueOf(data).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)

		label, skip := getTag(elem, field)
		if skip {
			continue
		}
		switch i {
		case 0:
			assert.Equal(t, "firstField", label)
		case 1:
			assert.Equal(t, "secondField", label)
		case 2:
			assert.Equal(t, "ThirdField", label)
		case 3:
			assert.Equal(t, "fourthField", label)
		}
	}

}
func Test_firstLowerCase(t *testing.T) {
	cases := []string{"firstField", "SecondField", "THIRDFIELD"}
	for i := range cases {
		switch i {
		case 0:
			assert.Equal(t, "firstField", firstLowerCase(cases[i]))
		case 1:
			assert.Equal(t, "secondField", firstLowerCase(cases[i]))
		case 2:
			assert.Equal(t, "tHIRDFIELD", firstLowerCase(cases[i]))
		}
	}
}
