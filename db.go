package sutando

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insertOneResult *mongo.InsertOneResult
type insertManyResult *mongo.InsertManyResult
type updateResult *mongo.UpdateResult
type deleteResult *mongo.DeleteResult

var (
	ErrNoDocument = mongo.ErrNoDocuments
)

type DB interface {
	/*
		Return mongo client diver
	*/
	GetDriver() *mongo.Client
	/*
		Return mongo database diver
	*/
	GetDriverDB() *mongo.Database
	/*
		Collection you want to operate.
	*/
	Collection(name string, opts ...*options.CollectionOptions) builder
	/*
		Insert data in MongoDB
			# Example:
				insert := db.Collection("CollectionName").Insert(&obj)
				_, _, err := db.ExecInsert(ctx, insert)

				insertMany := db.Collection("CollectionName").Insert(&obj1, &obj2, &obj3)
				_, _, err := db.ExecInsert(ctx, insertMany)
	*/
	ExecInsert(ctx context.Context, i *insert) (insertOneResult, insertManyResult, error)
	/*
		Find data in MongoDB
			# Example:
				result := struct{}
				query := db.Collection("CollectionName").Find().Equal("Name", "sutando").Greater("Number", 300)

				err := db.ExecFind(ctx, query, &result)
	*/
	ExecFind(ctx context.Context, q query, p any) error
	/*
		Update data in MongoDB
			# Example:
				update := db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello")
				result, err := su.db.ExecUpdate(su.ctx, update, false)
	*/
	ExecUpdate(ctx context.Context, u update, upsert bool) (updateResult, error)
	/*
		Delete data in MongoDB
			# Example:
				delete := db.Collection("Collection").Delete().Equal("Field", "sutando")
				result, err := su.db.ExecDelete(su.ctx, delete)
	*/
	ExecDelete(ctx context.Context, q query) (deleteResult, error)
	/*
		Disconnect closes sockets to the topology referenced by this Client. It will
		shut down any monitoring goroutines, close the idle connection pool, and will
		wait until all the in use connections have been returned to the connection
		pool and closed before returning. If the context expires via cancellation,
		deadline, or timeout before the in use connections have returned, the in use
		connections will be closed, resulting in the failure of any in flight read
		or write operations. If this method returns with no errors, all connections
		associated with this Client have been closed.
	*/
	Disconnect(ctx context.Context) error
}

type sutandoDB struct {
	client *mongo.Client
	db     string
}

/*
Create a new mongoDB connection

	# Sample Code:
		db, err := sutando.NewDB(ctx, sutando.Conn{
			Username:  "example",
			Password:  "example",
			Host:      "localhost",
			Port:      27017,	// leave empty if there's port in host
			DB:        "example",
			AdminAuth: true,
			Pem:       "",		// optional
		})

	# --- How To Use ---

	# Find:
		result := struct{}
		query := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First()
		err := db.ExecFind(ctx, query, &result)

	# Insert
		insert := db.Collection("Collection").Insert(&obj)
		result, _, err := db.ExecInsert(ctx, insert)

		insertMany := db.Collection("Collection").Insert(&obj1, &obj2, &obj3)
		_, resultMany, err := db.ExecInsert(ctx, insertMany)
	# Update with Model
		update := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").First()
		result, err := db.ExecUpdate(su.ctx, updateOne, false)

		updateMany := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando")
		result, err := db.ExecUpdate(su.ctx, updateMany, false)
	# Update with Set
		update := db.Collection("Collection").Update().Equal("Field", "sutando").First().Set("Field", "hello")
		result, err := db.ExecUpdate(su.ctx, updateOne, false)

		updateMany := db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello")
		result, err := db.ExecUpdate(su.ctx, updateMany, false)
	# Delete
		delete := db.Collection("Collection").Delete().Equal("Field", "sutando").First()
		result, err := db.ExecDelete(su.ctx, delete)

		deleteMany := db.Collection("Collection").Delete().Equal("Field", "sutando")
		result, err := db.ExecDelete(su.ctx, deleteMany)
*/
func NewDB(ctx context.Context, c Connection) (DB, error) {
	cfg := &tls.Config{
		RootCAs: x509.NewCertPool(),
	}
	var (
		opts *options.ClientOptions
		err  error
	)
	dsn, pem := c.DSN(cfg)

	opts = options.Client().ApplyURI(dsn).
		SetRegistry((*bsoncodec.Registry)(bson.NewRegistryBuilder().
			RegisterTypeEncoder(_TYPE_DECIMAL, coder{}).
			RegisterTypeDecoder(_TYPE_DECIMAL, coder{}).
			Build()))

	if pem {
		opts.SetTLSConfig(cfg)
	}

	c.SetupOption(opts)

	f := false
	opts.RetryWrites = &f
	client, err := mongo.Connect(
		ctx,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect, %w", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db.Ping, %w", err)
	}

	return &sutandoDB{client, c.Database()}, nil
}

/*
Create a mongoDB connection with an existed mongo-driver

	# --- How To Use ---

	# Sample Code:
		var client *mongo.Client
		...
		database := "example"
		db := sutando.NewDBFromMongo(ctx, client, database)

	# Find:
		result := struct{}
		query := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First()
		err := db.ExecFind(ctx, query, &result)

	# Insert
		insert := db.Collection("Collection").Insert(&obj)
		result, _, err := db.ExecInsert(ctx, insert)

		insertMany := db.Collection("Collection").Insert(&obj1, &obj2, &obj3)
		_, resultMany, err := db.ExecInsert(ctx, insertMany)
	# Update with Model
		update := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").First()
		result, err := db.ExecUpdate(su.ctx, updateOne, false)

		updateMany := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando")
		result, err := db.ExecUpdate(su.ctx, updateMany, false)
	# Update with Set
		update := db.Collection("Collection").Update().Equal("Field", "sutando").First().Set("Field", "hello")
		result, err := db.ExecUpdate(su.ctx, updateOne, false)

		updateMany := db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello")
		result, err := db.ExecUpdate(su.ctx, updateMany, false)
	# Delete
		delete := db.Collection("Collection").Delete().Equal("Field", "sutando").First()
		result, err := db.ExecDelete(su.ctx, delete)

		deleteMany := db.Collection("Collection").Delete().Equal("Field", "sutando")
		result, err := db.ExecDelete(su.ctx, deleteMany)
*/
func NewDBFromMongo(ctx context.Context, client *mongo.Client, database string) DB {
	return &sutandoDB{client, database}
}

func (s *sutandoDB) GetDriver() *mongo.Client {
	return s.client
}

func (s *sutandoDB) GetDriverDB() *mongo.Database {
	return s.client.Database(s.db)
}

func (s *sutandoDB) Collection(name string, opts ...*options.CollectionOptions) builder {
	return builder{col: s.client.Database(s.db).Collection(name, opts...)}
}

func (s *sutandoDB) ExecInsert(ctx context.Context, i *insert) (insertOneResult, insertManyResult, error) {
	var (
		err  error
		one  insertOneResult
		many insertManyResult
	)
	objects := i.build()
	switch len(objects) {
	case 0:
		return one, many, errors.New("object to insert should be pointer")
	case 1:
		one, err = i.col.InsertOne(ctx, objects[0], i.optionOne()...)
		return one, many, err
	default:
		many, err = i.col.InsertMany(ctx, objects, i.optionMany()...)
		return one, many, err
	}
}

func (s *sutandoDB) ExecFind(ctx context.Context, q query, p any) error {
	if reflect.TypeOf(p).Kind() != reflect.Pointer {
		return errors.New("object to find should be a pointer")
	}
	kind := reflect.TypeOf(p).Elem().Kind()
	if kind == reflect.Array {
		return errors.New("object to find cannot be an array")
	}

	if kind != reflect.Slice {
		if !q.isOne() {
			return errors.New("find too many results! use query with 'First' to find one result")
		}
		return execFindOne(ctx, q, p)
	}

	if !q.isOne() {
		return execFindMany(ctx, q, p)
	}

	obj := reflect.New(reflect.TypeOf(p).Elem().Elem())
	err := execFindOne(ctx, q, obj.Interface())
	if err != nil {
		return err
	}
	sli := reflect.Append(reflect.ValueOf(p).Elem(), obj.Elem())
	reflect.ValueOf(p).Elem().Set(sli)
	return nil
}

func execFindOne(ctx context.Context, q query, p any) error {
	result := q.col().FindOne(ctx, q.build())
	err := result.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

func execFindMany(ctx context.Context, q query, p any) error {
	cursor, err := q.col().Find(ctx, q.build())
	if err != nil {
		return err
	}
	err = cursor.All(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *sutandoDB) ExecUpdate(ctx context.Context, u update, upsert bool) (updateResult, error) {
	if u.isOne() {
		return u.col().UpdateOne(ctx, u.build(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	}
	return u.col().UpdateMany(ctx, u.build(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
}

func (s *sutandoDB) ExecDelete(ctx context.Context, q query) (deleteResult, error) {
	if q.isOne() {
		return q.col().DeleteOne(ctx, q.build(), &options.DeleteOptions{})
	}
	return q.col().DeleteMany(ctx, q.build(), &options.DeleteOptions{})
}

func (s *sutandoDB) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
