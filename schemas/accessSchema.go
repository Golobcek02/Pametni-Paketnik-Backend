package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Access struct {
	OwnerId   primitive.ObjectID `bson:"_id,omitempty"`
	AccessIds string
	BoxId     int
}
