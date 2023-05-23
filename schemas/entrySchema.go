package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entry struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DeliveryId   int
	BoxId        int
	Latitude     float64
	Longitude    float64
	TimeAccessed int64
	LoggerId     primitive.ObjectID
	EntryType    string
}
