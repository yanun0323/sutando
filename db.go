package sutando

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insertOneResult *mongo.InsertOneResult
type insertManyResult *mongo.InsertManyResult
type updateResult *mongo.UpdateResult

var (
	ErrNoDocument = mongo.ErrNoDocuments
)

type DB interface {
	/*
		Return mongo client diver
	*/
	Client() *mongo.Client
	/*
		Return mongo database diver
	*/
	DB() *mongo.Database
	/*
		The collection you want to operate.
	*/
	Collection(name string, opts ...*options.CollectionOptions) builder
	/*
		Insert data in MongoDB

			# Example:
				insert := e.Collection("CollectionName").Insert(&obj)
				_, _, err := e.ExecInsert(ctx, insert)

				insertMany := e.Collection("CollectionName").Insert(&obj1, &obj2, &obj3)
				_, _, err := e.ExecInsert(ctx, insertMany)
	*/
	ExecInsert(ctx context.Context, i *insert) (insertOneResult, insertManyResult, error)
	/*
		Find data in MongoDB

			# Example:
				result := struct{}
				find := e.Collection("CollectionName").Find().Equal("field_a", "sutando").Greater("field_b", 300)

				err := e.ExecFind(ctx, find, &result)
	*/
	ExecFind(ctx context.Context, q *find, p any) error
	ExecUpdate(ctx context.Context, u *update, upsert bool) (updateResult, error)
}

type sutandoDB struct {
	client *mongo.Client
	db     string
}

/*
Create a new connection to MongoDB

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
		find := e.Collection("CollectionName").Find().Equal("field_a", "sutando").Greater("field_b", 300)

		err := e.ExecFind(ctx, find, &result)

	# Insert
		insert := e.Collection("CollectionName").Insert(&obj)
		_, _, err := e.ExecInsert(ctx, insert)

		insertMany := e.Collection("CollectionName").Insert(&obj1, &obj2, &obj3)
		_, _, err := e.ExecInsert(ctx, insertMany)
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

func (s *sutandoDB) Client() *mongo.Client {
	return s.client
}

func (s *sutandoDB) DB() *mongo.Database {
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

	switch len(i.data) {
	case 0:
		return one, many, errors.New("object to insert should be pointer")
	case 1:
		one, err = i.col.InsertOne(ctx, i.build()[0], i.optionOne()...)
		return one, many, err
	default:
		many, err = i.col.InsertMany(ctx, i.build(), i.optionMany()...)
		return one, many, err
	}
}

func (s *sutandoDB) ExecFind(ctx context.Context, f *find, p any) error {
	if reflect.TypeOf(p).Kind() != reflect.Pointer {
		return errors.New("object to find should be a pointer")
	}
	kind := reflect.TypeOf(p).Elem().Kind()
	if kind == reflect.Array {
		return errors.New("object to find cannot be an array")
	}
	if kind == reflect.Slice {
		return execFindMany(ctx, f, p)
	}
	return execFindOne(ctx, f, p)

}

func execFindOne(ctx context.Context, q *find, p any) error {
	result := q.col.FindOne(ctx, q.build())
	err := result.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

func execFindMany(ctx context.Context, q *find, p any) error {
	cursor, err := q.col.Find(ctx, q.build())
	if err != nil {
		return err
	}
	err = cursor.All(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *sutandoDB) ExecUpdate(ctx context.Context, u *update, upsert bool) (updateResult, error) {
	switch len(u.data) {
	case 0:
		return nil, errors.New("object to insert should be pointer")
	case 1:
		if u.id > 0 {
			return u.col.UpdateByID(ctx, u.buildFindOrID(), u.buildObjects())
		}
		return u.col.UpdateOne(ctx, u.buildFindOrID(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	default:
		return u.col.UpdateMany(ctx, u.buildFindOrID(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	}
}
