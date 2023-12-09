package sutando

import "go.mongodb.org/mongo-driver/bson"

func bsonD(elems ...bson.E) bson.D {
	return append(bson.D{}, elems...)
}

func bsonM(key string, val any) bson.M {
	return bson.M{key: val}
}
