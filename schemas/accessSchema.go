package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Access struct {
	OwnerId   primitive.ObjectID
	AccessIds []primitive.ObjectID
	BoxId     int
}
