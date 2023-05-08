package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string
	Surname   string
	Username  string
	Email     string
	Password  string
	UserBoxes string
}
