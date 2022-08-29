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
	Collection(name string, opts ...*options.CollectionOptions) builder
	ExecInsert(ctx context.Context, i *insert) (insertOneResult, insertManyResult, error)
	ExecQuery(ctx context.Context, q *query, p any) error
	ExecUpdate(ctx context.Context, u *update, upsert bool) (updateResult, error)
}

type sutandoDB struct {
	client *mongo.Client
	db     string
}

/* sutando provide default connection implement 'Conn' */
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

/*
The collection you want to operate.
*/
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

func (s *sutandoDB) ExecQuery(ctx context.Context, q *query, p any) error {
	if reflect.TypeOf(p).Kind() != reflect.Pointer {
		return errors.New("object to find should be a pointer")
	}
	kind := reflect.TypeOf(p).Elem().Kind()
	if kind == reflect.Array {
		return errors.New("object to find cannot be an array")
	}
	if kind == reflect.Slice {
		return execQueryMany(ctx, q, p)
	}
	return execQueryOne(ctx, q, p)

}

func execQueryOne(ctx context.Context, q *query, p any) error {
	result := q.col.FindOne(ctx, q.build())
	err := result.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

func execQueryMany(ctx context.Context, q *query, p any) error {
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
			return u.col.UpdateByID(ctx, u.buildQueryOrID(), u.buildObjects())
		}
		return u.col.UpdateOne(ctx, u.buildQueryOrID(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	default:
		return u.col.UpdateMany(ctx, u.buildQueryOrID(), u.buildObjects(), &options.UpdateOptions{Upsert: &upsert})
	}
}
