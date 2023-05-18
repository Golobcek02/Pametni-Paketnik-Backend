package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	BoxID        int
	Status       string
	PackageRoute PackageRoutes
}
