package sutando

import "go.mongodb.org/mongo-driver/bson"

func bsonD(elems ...bson.E) bson.D {
	return append(bson.D{}, elems...)
}

func bsonE(key string, val any) bson.E {
	return bson.E{Key: key, Value: val}
}

func bsonM(key string, val any) bson.M {
	return bson.M{key: val}
}
