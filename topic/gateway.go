package topic

import (
	"NotifyMe-Bot/types"
	"github.com/pkg/errors"
)

var AlreadySubError = errors.New("Already subscribed. If you want to unsub from a topic use /unsubscribe")
var NoSubError = errors.New("You weren't subscribed to that topic!")

func GetAlertMessage(tag string, groupID int64) string {
	topicToAlert, err := FindByTopic(tag, groupID)
	if err != nil {
		return ""
	}
	return prepareAlertMsg(topicToAlert)
}

func CreateTopic(topicName string, groupID int64, user types.User) error {

	newTopic := Topic{
		Name:        topicName,
		Creator:     user.UserName,
		Subscribers: []string{user.UserName},
		TimesCalled: 0,
		GroupID:     groupID,
	}

	_,err := Save(newTopic)
	if err != nil{
		return err
	}


	return nil
}


func SubscribeToTopic(topicName string, groupID int64, user types.User) error {
	topic,err := FindByTopic(topicName,groupID)
	if err != nil{
		return err
	}
	if subscribed(topic.Subscribers, user.UserName){
		return AlreadySubError
	}
	topic.Subscribers = append(topic.Subscribers, user.UserName)
	err = Update(topic)
	if err != nil{
		return err
	}
	return nil
}

func UnsubscribeToTopic(topicName string, groupID int64, user types.User) error {
	topic,err := FindByTopic(topicName,groupID)
	if err != nil{
		return err
	}
	if !subscribed(topic.Subscribers, user.UserName){
		return NoSubError
	}

	for i,_ := range topic.Subscribers{
		if topic.Subscribers[i] == user.UserName{
			topic.Subscribers = append(topic.Subscribers[:i], topic.Subscribers[i+1:]...)
			break
		}
	}
	err = Update(topic)
	if err != nil{
		return err
	}
	return nil
}

func subscribed(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func prepareAlertMsg(topic Topic) string {
	var msg string
	for _,sub := range topic.Subscribers{
		msg = msg + "@" + sub + " "
	}

	return msg
}