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
	FifthField  string `bson:"-"`
}

func TestGetTag(t *testing.T) {
	data := &testTagStruct{
		FirstField:  "1",
		SecondField: "2",
		thirdField:  "3",
		fourthField: "4",
		FifthField:  "5",
	}
	elem := reflect.ValueOf(data).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)

		label, skip, _, _ := getTag(elem, field)
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
func TestFirstLowerCase(t *testing.T) {
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
