package sutando

import (
	"context"

	"github.com/yanun0323/sutando/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type inserting interface {
	// Exec executes inserting documents
	// and return corresponding result of inserting.
	Exec(context.Context) (insertOneResult, insertManyResult, error)
}

type updating interface {
	// Exists matches documents that contain or
	// do not contain a specified field,
	// including documents where the field value is null.
	Exists(key string, exists bool) updating

	// And matches documents where the value of a field
	// equals the specified value.
	And(key string, value any) updating

	// Equal matches documents where the value of a field
	// equals the specified value.
	Equal(key string, value any) updating

	// NotEqual selects the documents where the value
	// of the specified field is not equal to the specified value.
	// This includes documents that do not contain the specified field.
	NotEqual(key string, value ...any) updating

	// Greater selects those documents where the value
	// of the specified field is greater than
	// (i.e. >) the specified value.
	Greater(key string, value any) updating

	// GreaterOrEqual selects the documents where the value
	// of the specified field is greater than or equal to
	// (i.e. >=) a specified value (e.g. value.)
	GreaterOrEqual(key string, value any) updating

	// Less selects the documents where the value of the field is less
	// than (i.e. <) the specified value.
	Less(key string, value any) updating

	// LessOrEqual selects the documents where the value of the field
	// is less than or equal to (i.e. <=) the specified value.
	LessOrEqual(key string, value any) updating

	// Contain selects the documents where the value of a field
	// is an array that contains all the specified elements.
	Contain(key string, value ...any) updating

	// In selects the documents where the value of a field
	// equals any value in the specified array.
	In(key string, value ...any) updating

	// NotIn selects the documents where:
	//
	// - the specified field value is not in the specified array or
	//
	// - the specified field does not exist.
	NotIn(key string, value ...any) updating

	// Regex provides regular expression capabilities
	// for pattern matching strings in queries.
	Regex(key string, regex string, opt ...option.Regex) updating

	// First selects the first document.
	First() updating

	// Set sets keys to update with input values.
	//
	// (no effect when the function start with UpdateWith())
	Set(key string, value any) updating

	// Exec executes updating documents
	// and return the result of updating.
	Exec(ctx context.Context, upsert bool) (updateResult, error)

	// buildObjects encode the objects with sutando encoder.
	buildObjects() any
}

type querying interface {
	// col return the collection driver from the connection.
	col() *mongo.Collection

	// build builds the filter structure with previous operations.
	build() bson.D

	// isOne return whether the 'First' method called.
	isOne() bool

	// Exists matches documents that contain or
	// do not contain a specified field,
	// including documents where the field value is null.
	Exists(key string, exists bool) querying

	// And matches documents where the value of a field
	// equals the specified value.
	And(key string, value any) querying

	// Equal matches documents where the value of a field
	// equals the specified value.
	Equal(key string, value any) querying

	// NotEqual selects the documents where the value
	// of the specified field is not equal to the specified value.
	// This includes documents that do not contain the specified field.
	NotEqual(key string, value ...any) querying

	// Greater selects those documents where the value
	// of the specified field is greater than
	// (i.e. >) the specified value.
	Greater(key string, value any) querying

	// GreaterOrEqual selects the documents where the value
	// of the specified field is greater than or equal to
	// (i.e. >=) a specified value (e.g. value.)
	GreaterOrEqual(key string, value any) querying

	// Less selects the documents where the value of the field is less
	// than (i.e. <) the specified value.
	Less(key string, value any) querying

	// LessOrEqual selects the documents where the value of the field
	// is less than or equal to (i.e. <=) the specified value.
	LessOrEqual(key string, value any) querying

	// Contain selects the documents where the value of a field
	// is an array that contains all the specified elements.
	Contain(key string, value ...any) querying

	// In selects the documents where the value of a field
	// equals any value in the specified array.
	In(key string, value ...any) querying

	// NotIn selects the documents where:
	//
	// - the specified field value is not in the specified array or
	//
	// - the specified field does not exist.
	NotIn(key string, value ...any) querying

	// Regex provides regular expression capabilities
	// for pattern matching strings in queries.
	Regex(key string, regex string, opt ...option.Regex) querying

	// Count counts documents and return the result.
	Count(ctx context.Context, index ...string) (int64, error)

	// First selects the first document.
	first() querying
}

type finding interface {
	// Exists matches documents that contain or
	// do not contain a specified field,
	// including documents where the field value is null.
	Exists(key string, exists bool) finding

	// And matches documents where the value of a field
	// equals the specified value.
	And(key string, value any) finding

	// Equal matches documents where the value of a field
	// equals the specified value.
	Equal(key string, value any) finding

	// NotEqual selects the documents where the value
	// of the specified field is not equal to the specified value.
	// This includes documents that do not contain the specified field.
	NotEqual(key string, value ...any) finding

	// Greater selects those documents where the value
	// of the specified field is greater than
	// (i.e. >) the specified value.
	Greater(key string, value any) finding

	// GreaterOrEqual selects the documents where the value
	// of the specified field is greater than or equal to
	// (i.e. >=) a specified value (e.g. value.)
	GreaterOrEqual(key string, value any) finding

	// Less selects the documents where the value of the field is less
	// than (i.e. <) the specified value.
	Less(key string, value any) finding

	// LessOrEqual selects the documents where the value of the field
	// is less than or equal to (i.e. <=) the specified value.
	LessOrEqual(key string, value any) finding

	// Contain selects the documents where the value of a field
	// is an array that contains all the specified elements.
	Contain(key string, value ...any) finding

	// In selects the documents where the value of a field
	// equals any value in the specified array.
	In(key string, value ...any) finding

	// NotIn selects the documents where:
	//
	// - the specified field value is not in the specified array or
	//
	// - the specified field does not exist.
	NotIn(key string, value ...any) finding

	// Regex provides regular expression capabilities
	// for pattern matching strings in queries.
	Regex(key string, regex string, opt ...option.Regex) finding

	// First selects the first document.
	First() finding

	// Sort sorts documents with input sorting parameters.
	//
	//	sort := map[string]bool{"name": true, "age": false} /* name ascend, age descend */
	//	db.Find().Sort(sort)
	Sort(asc map[string]bool) finding

	// Limit limits the number of documents.
	Limit(i int64) finding

	// Skip skips the number of documents.
	Skip(i int64) finding

	// Exec executes finding documents and parsing the result with input pointer structure.
	Exec(ctx context.Context, result any) error
}

type deleting interface {
	// Exists matches documents that contain or
	// do not contain a specified field,
	// including documents where the field value is null.
	Exists(key string, exists bool) deleting

	// equals the specified value.
	// equals the specified value.
	And(key string, value any) deleting

	// Equal matches documents where the value of a field
	// equals the specified value.
	Equal(key string, value any) deleting

	// NotEqual selects the documents where the value
	// of the specified field is not equal to the specified value.
	// This includes documents that do not contain the specified field.
	NotEqual(key string, value ...any) deleting

	// Greater selects those documents where the value
	// of the specified field is greater than
	// (i.e. >) the specified value.
	Greater(key string, value any) deleting

	// GreaterOrEqual selects the documents where the value
	// of the specified field is greater than or equal to
	// (i.e. >=) a specified value (e.g. value.)
	GreaterOrEqual(key string, value any) deleting

	// Less selects the documents where the value of the field is less
	// than (i.e. <) the specified value.
	Less(key string, value any) deleting

	// LessOrEqual selects the documents where the value of the field
	// is less than or equal to (i.e. <=) the specified value.
	LessOrEqual(key string, value any) deleting

	// Contain selects the documents where the value of a field
	// is an array that contains all the specified elements.
	Contain(key string, value ...any) deleting

	// In selects the documents where the value of a field
	// equals any value in the specified array.
	In(key string, value ...any) deleting

	// NotIn selects the documents where:
	//
	// - the specified field value is not in the specified array or
	//
	// - the specified field does not exist.
	NotIn(key string, value ...any) deleting

	// Regex provides regular expression capabilities
	// for pattern matching strings in queries.
	Regex(key string, regex string, opt ...option.Regex) deleting

	// First selects the first document.
	First() deleting

	// Exec executes deleting documents
	Exec(context.Context) (deleteResult, error)
}
