package sutando

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
		Deprecated after sutando v1.3.0, using code below to instead:
			_, _, err := db.Collection("collectionName").Insert(&data).Exec(ctx)
	*/
	ExecInsert(ctx context.Context, i inserting) (insertOneResult, insertManyResult, error)

	/*
		Deprecated after sutando v1.3.0, using code below to instead:
			err := db.Collection("collectionName").Find().First().Exec(ctx, &result)
	*/
	ExecFind(ctx context.Context, f finding, p any) error

	/*
		Deprecated after sutando v1.3.0, using code below to instead:
			_, err := db.Collection("collectionName").Update().Equal("Field", "sutando").Set("Field", "hello").First().Exec(ctx, false)
	*/
	ExecUpdate(ctx context.Context, u updating, upsert bool) (updateResult, error)

	/*
		Deprecated after sutando v1.3.0, using code below to instead:
			_, err := db.Collection("collectionName").Delete().First().Exec(ctx)
	*/
	ExecDelete(ctx context.Context, d deleting) (deleteResult, error)
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

	// Sample: Using Host and Port
	db, err := sutando.NewDB(ctx, sutando.Conn{
		Username:  "example",
		Password:  "example",
		Host:      "example",
		Port:      27017,
		DB:        "example",
		AdminAuth: true,
		Pem:       "",		// optional
		ClientOptionsHandler: func(opts *options.ClientOptions) {
			// do something...
		},
	})

	// Sample: Using SRV URL
	db, err := sutando.NewDB(ctx, sutando.Conn{
		Username:  "example",
		Password:  "example",
		Host:      "example.mongo.net",
		DB:        "example",
		Srv:       true,
		ClientOptionsHandler: func(opts *options.ClientOptions) {
			// do something...
		},
	})
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

	reg := bson.NewRegistry()
	reg.RegisterTypeEncoder(_TYPE_DECIMAL, coder{})
	reg.RegisterTypeDecoder(_TYPE_DECIMAL, coder{})

	opts = options.Client().ApplyURI(dsn).
		SetRegistry(reg)

	if pem {
		opts.SetTLSConfig(cfg)
	}

	c.SetupClientOptions(opts)

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

	// Sample:
	var client *mongo.Client
	...
	database := "example"
	db := sutando.NewDBFromMongo(ctx, client, database)
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

func (s *sutandoDB) ExecInsert(ctx context.Context, i inserting) (insertOneResult, insertManyResult, error) {
	return i.Exec(ctx)
}

func (s *sutandoDB) ExecFind(ctx context.Context, f finding, p any) error {
	return f.Exec(ctx, p)
}

func (s *sutandoDB) ExecUpdate(ctx context.Context, u updating, upsert bool) (updateResult, error) {
	return u.Exec(ctx, upsert)
}

func (s *sutandoDB) ExecDelete(ctx context.Context, d deleting) (deleteResult, error) {
	return d.Exec(ctx)
}

func (s *sutandoDB) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
