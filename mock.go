package sutando

import (
	"time"

	"github.com/shopspring/decimal"
)

type testStruct struct {
	Name      string    `bson:"nameName"`
	Age       int       `bson:"age"`
	Birthday  time.Time `bson:"birthday"`
	Ignore    int       `bson:"-"`
	Inner     testSubStruct
	Inner2    testSubStruct
	Arr       []int
	Map       map[string]int
	StructMap map[int]testSubStruct
	FitValue  decimal.Decimal
}

type testSubStruct struct {
	Name string `bson:"name"`
	Age  int
}

func mockData() testStruct {
	return testStruct{
		Name:     "Yanun",
		Age:      27,
		Birthday: time.Date(1995, time.March, 23, 0, 0, 0, 0, time.Local),
		Ignore:   10,
		Inner: testSubStruct{
			Name: "inner",
			Age:  50,
		},
		Inner2: testSubStruct{
			Name: "inner2",
			Age:  10,
		},
		Arr: []int{1, 2, 3, 4, 5},
		Map: map[string]int{"1": 2, "3": 4, "5": 6},
		StructMap: map[int]testSubStruct{
			1: {
				Name: "inner",
				Age:  50,
			},
			2: {
				Name: "inner2",
				Age:  10,
			},
		},
		FitValue: decimal.RequireFromString("0.03"),
	}
}
