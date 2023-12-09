package sutando

import "go.mongodb.org/mongo-driver/bson"

type filter struct {
	d bson.D
}

func (f *filter) append(key string, val any) {
	f.d = append(f.d, bson.E{Key: key, Value: val})
}
