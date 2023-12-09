package sutando

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insertOneResult *mongo.InsertOneResult
type insertManyResult *mongo.InsertManyResult
type updateResult *mongo.UpdateResult
type deleteResult *mongo.DeleteResult

var (
	// ErrNoDocuments occurs when the operation
	// that created the SingleResult did not return
	// any document.
	ErrNoDocument = mongo.ErrNoDocuments

	// ErrDuplicatedKey occurs when there is a
	// unique key constraint violation.
	ErrDuplicatedKey = errors.New("mongo: duplicated key")
)

// DB provides a definition of db behavior.
type DB interface {
	// RawClient returns the raw mongo client diver.
	RawClient() *mongo.Client

	// RawDatabase returns the raw mongo database diver.
	RawDatabase() *mongo.Database

	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
	//
	//	col := db.Collection("col_name")
	//
	//	// insert
	//		col.Insert(&obj).Exec(ctx)
	//
	//	// find
	//		col.Find().Equal("name", "yanun").First().Exec(ctx, &result)
	//
	//	// update
	//		upsert := true
	//		col.Update().Set("name", "changed").Exec(ctx, upsert)
	//		col.UpdateWith(&obj).Exec(ctx, upsert)
	//
	//	// delete
	//		col.Delete().Greater("age", 20).Exec(ctx)
	Collection(name string, opts ...*options.CollectionOptions) builder

	// Disconnect closes sockets to the topology referenced by this Client. It will
	// shut down any monitoring goroutines, close the idle connection pool, and will
	// wait until all the in use connections have been returned to the connection
	// pool and closed before returning. If the context expires via cancellation,
	// deadline, or timeout before the in use connections have returned, the in use
	// connections will be closed, resulting in the failure of any in flight read
	// or write operations. If this method returns with no errors, all connections
	// associated with this Client have been closed.
	Disconnect(ctx context.Context) error
}

type sutandoDB struct {
	client *mongo.Client
	db     string
}

// NewDB initializes a sutando.DB from providing connection configuration.
//
//	// connect through host and port.
//	db, err := sutando.NewDB(ctx, sutando.Conn{
//		Username:  "example",
//		Password:  "example",
//		Host:      "example",
//		Port:      27017,
//		DB:        "example",
//		AdminAuth: true,
//		Pem:       "",
//		ClientOptionsHandler: func(opts *options.ClientOptions) {
//			opts.SetConnectTimeout(5 * time.Second)
//			opts.SetTimeout(15 * time.Second)
//		},
//	})
//
//	// connect through SRV.
//	db, err := sutando.NewDB(ctx, sutando.ConnSrv{
//		Username:  "example",
//		Password:  "example",
//		Host:      "example.mongo.net",
//		DB:        "example",
//		AdminAuth: true,
//		Pem:       "",
//		ClientOptionsHandler: func(opts *options.ClientOptions) {
//			opts.SetConnectTimeout(5 * time.Second)
//			opts.SetTimeout(15 * time.Second)
//		},
//	})
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

// NewDBFromMongo initializes a sutando.DB by using the mongoDB connection from an existed mongo-driver.
//
//	var client *mongo.Client
//	...
//	database := "example"
//	db := sutando.NewDBFromMongo(ctx, client, database)
func NewDBFromMongo(ctx context.Context, client *mongo.Client, database string) DB {
	return &sutandoDB{client, database}
}

func (s *sutandoDB) RawClient() *mongo.Client {
	return s.client
}

func (s *sutandoDB) RawDatabase() *mongo.Database {
	return s.client.Database(s.db)
}

func (s *sutandoDB) Collection(name string, opts ...*options.CollectionOptions) builder {
	return builder{col: s.client.Database(s.db).Collection(name, opts...)}
}

func (s *sutandoDB) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
