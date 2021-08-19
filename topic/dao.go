package topic

import (
	"NotifyMe-Bot/client/db"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	collection = "topic"
)

type Dao interface{
	Save(topic Topic)(string,error)
	FindByCreator(creator string,groupID int64)(Topic,error)
	FindByTopic(topicName string,groupID int64)(Topic,error)

}

func Save(topic Topic) (string,error){
	collection := db.GetCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, topic)
	if err != nil {
		return "",err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if  !ok {
		return "",errors.New("Could not getInsertedID")
	}
	return oid.String(),nil

}

func Update(topic Topic) error{
	collection := db.GetCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id":topic.Id},
		bson.D{{"$set",topic}})
	if err != nil {
		return err
	}



	return nil
}

func FindByCreator(creator string,groupID int64) (Topic,error) {
	collection := db.GetCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topic := &Topic{}
	filter := bson.D{{"creator",creator},{"group_id",groupID}}

	err := collection.FindOne(ctx, filter).Decode(&topic)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			return Topic{},errors.Wrap(err,"No documents found")
		}
		return Topic{},errors.Wrap(err,"could not connect to DB")
	}

	return *topic,nil
}

func FindByTopic(topicName string,groupID int64) (Topic,error) {
	collection := db.GetCollection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topic := &Topic{}
	filter := bson.D{{"name",topicName},{"group_id",groupID}}

	err := collection.FindOne(ctx, filter).Decode(&topic)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			return Topic{},errors.Wrap(err,"No documents found")
		}
		return Topic{},errors.Wrap(err,"could not connect to DB")
	}

	return *topic,nil
}