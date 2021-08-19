package topic

import "go.mongodb.org/mongo-driver/bson/primitive"

type Topic struct{
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Creator string `bson:"creator"`
	Subscribers []string `bson:"subscribers"`
	TimesCalled int `bson:"times_called"`
	GroupID int64 `bson:"group_id"`
}