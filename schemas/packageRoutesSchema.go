package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type PackageRoutes struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Stops  []string
	Orders []Order
}
