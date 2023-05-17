package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Box struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BoxId     int
	Latitude  float64
	Longitude float64
	OwnerId   primitive.ObjectID
	AccessIds []primitive.ObjectID
}
