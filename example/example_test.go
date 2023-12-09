package example

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/yanun0323/sutando"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testInstant(); err != nil {
			b.Fatal(err)
		}
	}
}

func Test(t *testing.T) {
	if err := testInstant(); err != nil {
		t.Fatal(err)
	}
}

func testInstant() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sutando.NewDB(ctx, &sutando.Conn{
		Username:  "test",
		Password:  "test",
		Host:      "localhost",
		Port:      27017,
		DB:        "sutando",
		AdminAuth: true,
		ClientOptionsHandler: func(opts *options.ClientOptions) {
			opts.SetConnectTimeout(5 * time.Second)
			opts.SetTimeout(15 * time.Second)
		},
	})
	if err != nil {
		return (err)
	}

	collection := strconv.FormatInt(time.Now().Unix(), 10)
	col := db.Collection(collection)
	_, _ = col.Delete().Exec(ctx)

	{ // Insert
		result, _, err := col.Insert(_school).Exec(ctx)
		if err != nil {
			return (err)
		}
		if result.InsertedID == nil {
			return errors.New("empty inserted ID")
		}
	}

	{ // Find
		var result School
		if err := col.Find().Equal("name", "sutando").Exists("room.901", true).First().Exec(ctx, &result); err != nil {
			return err
		}
	}

	{ // Update
		_school.Name = "changed"
		result, err := col.UpdateWith(&_school).Equal("name", "sutando").Exec(ctx, false)
		if err != nil {
			return err
		}

		if result.ModifiedCount != 1 {
			return errors.New("updated nothing")
		}

		var found School
		if err := col.Find().Equal("name", "sutando").First().Exec(ctx, &found); !errors.Is(sutando.ErrNoDocument, err) {
			return err
		}
	}

	{ // Delete
		result, err := col.Delete().Contain("room").Exec(ctx)
		if err != nil {
			return err
		}
		if result.DeletedCount != 1 {
			return errors.New("deleted nothing")
		}
	}

	return col.Drop(ctx)
}

var (
	_school = School{
		Name: "sutando",
		Room: map[string]*Class{
			"901": &_classGolang,
			"902": &_classSre,
		},
		OpenDay:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		BuildedTime: time.Now(),
		Money:       decimal.NewFromFloat(111_801_428.77062),
	}

	_classGolang = Class{
		MainTeacher:   &_teacherVic,
		OtherTeachers: []*Teacher{},
		Students:      []*Student{&_studentYanun, &_studentKai, &_studentVictor},
	}

	_classSre = Class{
		MainTeacher:   &_teacherMicheal,
		OtherTeachers: []*Teacher{&_teacherVic},
		Students:      []*Student{&_studentHarlan, &_studentTina},
	}

	_teacherMicheal = Teacher{
		People: People{
			Name: "Micheal",
			Age:  50,
		},
	}
	_teacherVic = Teacher{
		People: People{
			Name: "Vic",
			Age:  44,
		},
	}

	_studentHarlan = Student{
		People: People{
			Name: "Harlan",
			Age:  22,
		},
	}

	_studentTina = Student{
		People: People{
			Name: "Tina",
			Age:  19,
		},
	}

	_studentYanun = Student{
		People: People{
			Name: "Yanun",
			Age:  20,
		},
	}

	_studentKai = Student{
		People: People{
			Name: "Kai",
			Age:  25,
		},
	}

	_studentVictor = Student{
		People: People{
			Name: "Victor",
			Age:  21,
		},
	}
)

type School struct {
	Name        string
	Room        map[string]*Class
	OpenDay     []time.Weekday
	BuildedTime time.Time
	Money       decimal.Decimal
}

type Class struct {
	MainTeacher   *Teacher
	OtherTeachers []*Teacher `bson:",omitempty"`
	Students      []*Student `bson:",omitempty"`
}

type People struct {
	Name string
	Age  int
}

type Teacher struct {
	People
}

type Student struct {
	People
}
