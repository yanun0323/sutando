package sutando

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type builder struct {
	col *mongo.Collection
}

/*
Insert data
*/
func (b builder) Insert(p ...any) inserting {
	return newInsert(b.col, p...)
}

/*
Update data with model (Will update all fields including empty fields)
*/
func (b builder) UpdateWith(p any) updating {
	return newUpdate(b.col, p)
}

/*
Update data with set
*/
func (b builder) Update() updating {
	return newUpdate(b.col, nil)
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
