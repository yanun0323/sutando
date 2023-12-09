package example

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/yanun0323/sutando"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Benchmark(b *testing.B) {
	db := connect(b.Fatal)
	for i := 0; i < b.N; i++ {
		if err := testInstant(db, strconv.Itoa(rand.Int())); err != nil {
			b.Fatal(err)
		}
	}
}

func Test(t *testing.T) {
	t.Log("Test")
	db := connect(t.Fatal)
	if err := testInstant(db, "example_collection"); err != nil {
		t.Fatalf("%+v", err)
	}
}

var _ sutando.Connection = (*testStruct)(nil)

type testStruct struct {
}

func (ts *testStruct) DSN(cfg *tls.Config) (string, bool) /* string: dsn, bool: isPem */ {
	return "", true
}

func (ts *testStruct) Database() string {
	return ""
}

func (ts *testStruct) SetupClientOptions(*options.ClientOptions) {}

func connect(fatal func(args ...any)) sutando.DB {
	db, err := sutando.NewDB(context.Background(), &sutando.Conn{
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
		fatal(err)
	}
	return db
}

func testInstant(db sutando.DB, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := db.Collection(collection)
	defer col.Drop(ctx)

	_, _ = col.Delete().Exec(ctx)

	{ // Insert
		result, _, err := col.Insert(_school).Exec(ctx)
		if err != nil {
			return errors.New(fmt.Sprintf("%+v", err))
		}
		if result.InsertedID == nil {
			return errors.New("empty inserted ID")
		}
	}

	{ // Count
		c, err := col.Scalar().Count(ctx)
		if err != nil {
			return errors.New(fmt.Sprintf("%+v", err))
		}
		if c == 0 {
			errors.New("no data in database")
		}
	}

	{ // Find
		var result School
		err := col.Find().Equal("name", "sutando").Exists("room.901", true).First().Exec(ctx, &result)
		if err != nil {
			return errors.New(fmt.Sprintf("%+v", err))
		}
	}

	{ // Update
		_school.BuildedTime = time.Now()
		result, err := col.UpdateWith(&_school).Equal("name", "sutando").Exec(ctx, false)
		if err != nil {
			return errors.New(fmt.Sprintf("%+v", err))
		}

		if result.MatchedCount == 0 {
			return errors.New("updated nothing")
		}

		var found School
		err = col.Find().Equal("name", "sutando").First().Exec(ctx, &found)
		if err != nil && !errors.Is(sutando.ErrNoDocument, err) {
			return errors.New(fmt.Sprintf("%+v", err))
		}
	}

	{ // Delete
		result, err := col.Delete().Contain("room").Exec(ctx)
		if err != nil {
			return errors.New(fmt.Sprintf("%+v", err))
		}
		if result.DeletedCount != 1 {
			return errors.New("deleted nothing")
		}
	}

	return nil
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
