package telegram

import (
	"notify-me-bot/topic"
	"notify-me-bot/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

const newTopicHelp = "/newTopic TopicName"

func HandleUpdate(res http.ResponseWriter, req *http.Request){
	fullBody := &types.Update{}

	if err := json.NewDecoder(req.Body).Decode(fullBody); err != nil {
		fmt.Println("could not decode request full body", err)
		return
	}

	if fullBody.Message == nil {
		return
	}
	if strings.ToLower(fullBody.Message.Text) == "soy lu" {
		sendMessage("Hola, mi creador te ama mucho",fullBody.Message.Chat)
	}


	if strings.HasPrefix(fullBody.Message.Text,"#"){
		if len(strings.Split(fullBody.Message.Text, " ")) > 1{
			return
		}

		tag := strings.Trim(fullBody.Message.Text,"#")
		alertMsg := topic.GetAlertMessage(tag,fullBody.Message.Chat.ID)
		err := sendMessage(alertMsg,fullBody.Message.Chat)
		if err != nil {
			fmt.Println("error in sending reply:", err)
		}
	}

	if strings.HasPrefix(fullBody.Message.Text,"/"){
		command := strings.Trim(strings.Split(fullBody.Message.Text, " ")[0] , "/")
		args := strings.Split(fullBody.Message.Text, " ")[1:]

		if fullBody.Message.From.UserName == ""{
			sendMessage("It seems that you dont have an username created. You can't use this bot for now. Sorry!",fullBody.Message.Chat)
			return
		}
		switch strings.ToLower(command) {
		case "subscribe":
			if len(args) < 1{
				sendMessage("You need to provide a topic name! /subscribe MyTopicName" ,fullBody.Message.Chat)
				return
			}
			err := topic.SubscribeToTopic(args[0],fullBody.Message.Chat.ID,*fullBody.Message.From)
			if err != nil {
				if err == topic.AlreadySubError{
					sendMessage(err.Error(),fullBody.Message.Chat)
					return
				}
				sendMessage("error while subscribing to topic :( - error:" + err.Error(),fullBody.Message.Chat)
				return
			}
			sendMessage("Subscribed!",fullBody.Message.Chat)
			return
		case "unsubscribe":
			if len(args) < 1{
				sendMessage("You need to provide a topic name! /subscribe MyTopicName" ,fullBody.Message.Chat)
				return
			}
			err := topic.UnsubscribeToTopic(args[0],fullBody.Message.Chat.ID,*fullBody.Message.From)
			if err != nil {
				if err == topic.NoSubError{
					sendMessage(err.Error(),fullBody.Message.Chat)
					return
				}
				sendMessage("error while unsubscribing to topic :( - error:" + err.Error(),fullBody.Message.Chat)
				return
			}
			sendMessage("Unsubscribed!",fullBody.Message.Chat)
			return
		case "mysubs":
			print("3")
		case "newtopic":
			fmt.Println("InNewTopic")
			if args[0] == "help"{
				sendMessage(newTopicHelp,fullBody.Message.Chat)
				return
			}
			err := topic.CreateTopic(args[0],fullBody.Message.Chat.ID,*fullBody.Message.From)
			if err != nil{
				sendMessage("error while creating topic :( - error:" + err.Error(),fullBody.Message.Chat)
				return
			}
			sendMessage("Topic #%s was created!" + args[0],fullBody.Message.Chat)
			return

		default:
			return
		}
	}
}

type SimpleMessage struct {
	ChatID interface{} `json:"chat_id"`
	Text   string `json:"text"`
}
type MentionMsg struct {
	ChatID interface{} `json:"chat_id"`
	Text   string `json:"text"`
	Entities types.MessageEntity `json:"entities"`
}

func sendMessage(text string ,chat *types.Chat) error {
	// Create the request body struct

	reqBody := &SimpleMessage{
		ChatID: chat.ID,
		Text:   text,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1495263911:AAEjmxdCuPazMzeegGeZpm1RLBBMtFgx2oE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {

		return err
	}
	if res.StatusCode != http.StatusOK {
		var builder = new(strings.Builder)
		io.Copy(builder,res.Body)
		body := builder.String()
		fmt.Println(body)
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func sendMessageEntity(text string ,chat *types.Chat,user *types.User) error {
	// Create the request body struct

	reqBody := &MentionMsg{
		ChatID: chat.ID,
		Text:   text,
		Entities: types.MessageEntity{Type: "text_mention",User: user},
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1495263911:AAEjmxdCuPazMzeegGeZpm1RLBBMtFgx2oE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {

		return err
	}
	if res.StatusCode != http.StatusOK {
		var builder = new(strings.Builder)
		io.Copy(builder,res.Body)
		body := builder.String()
		fmt.Println(body)
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
