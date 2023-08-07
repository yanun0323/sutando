package sutando

import (
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testStruct struct {
	StructName   string
	StructAge    int
	RealBirthday time.Time
	Ignore       int `bson:"-"`
	primitive.M
	Inner     testSubStruct `bson:",inline"`
	ArrTest   []int
	MapTest   map[string]int
	StructMap map[int]testSubStruct
	FitValue  decimal.Decimal
	Random    int64
}

type testSubStruct struct {
	JustMyName string `bson:"name"`
	OhTheAge   int
}

func mockData() testStruct {
	return testStruct{
		StructName:   "Yanun",
		StructAge:    27,
		RealBirthday: time.Date(1995, time.March, 23, 0, 0, 0, 0, time.Local),
		Ignore:       10,
		M:            primitive.M{"Hello": 123},
		Inner: testSubStruct{
			JustMyName: "inner",
			OhTheAge:   10,
		},
		ArrTest: []int{1, 2, 3, 4, 5},
		MapTest: map[string]int{"1": 2, "3": 4, "5": 6},
		StructMap: map[int]testSubStruct{
			1: {
				JustMyName: "inner",
				OhTheAge:   50,
			},
			2: {
				JustMyName: "inner2",
				OhTheAge:   10,
			},
		},
		FitValue: decimal.RequireFromString("0.03"),
		Random:   time.Now().UnixMilli(),
	}
}
