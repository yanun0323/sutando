package sutando

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

/*
Set Custom Encode Types/Structure
*/
func (b builder) SetCustomEncodeTypes() *builder {
	return &b
}

/*
Insert data
*/
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, newBsonEncoder(), p...)
}

/*
Update data with model (Will update all fields including empty fields)
*/
func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, newBsonEncoder(), p)
}

/*
Update data with set
*/
func (b builder) Update() updating {
	return newUpdate(b.col, newBsonEncoder(), nil)
}

/*
Find data
*/
func (b builder) Find() finding {
	return newFind(b.col)
}

/*
Delete data
*/
func (b builder) Delete() deleting {
	return newDelete(b.col)
}

/*
For Query Test
*/
func (b builder) Query() querying {
	return newQuery(b.col)
}
