package main

import (
	"NotifyMe-Bot/client/db"
	"NotifyMe-Bot/telegram"
	"NotifyMe-Bot/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	fullBody := &types.Update{}

	if err := json.NewDecoder(req.Body).Decode(fullBody); err != nil {
		fmt.Println("could not decode request full body", err)
		return
	}
	fmt.Println("PRINTING REQUEST FULL BODY")
	fmt.Printf("%+v\n", fullBody)
	if fullBody.Message == nil {
		return
	}
	fmt.Printf("%+v\n", fullBody.Message.Text)
	// Check if the message contains the word "marco"
	// if not, return without doing anything
	if strings.Contains(fullBody.Message.Text,"testDB"){
		db.Add(fullBody.Message.Text)
	}
	if strings.Contains(fullBody.Message.Text,"testTag"){
		err := say(fullBody.Message.Chat,fullBody.Message.From, fmt.Sprintf("[inline mention of a user](tg://user?id=%d)",fullBody.Message.From.ID))
		if err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	}
	if !strings.Contains(strings.ToLower(fullBody.Message.Text), "marco") {
		return
	}

	// If the text contains marco, call the `sayPolo` function, which
	// is defined below
	if err := sayPolo(fullBody.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}


	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

//The below code deals with the process of sending a response message
// to the user

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// sayPolo takes a chatID and sends "polo" to them
func sayPolo(chatID int64) error {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   "Polo!!",
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
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
func say(chat *types.Chat,user *types.User,msg string) error {
	// Create the request body struct
	//reqBody := &sendMessageReqBody{
	//	ChatID: chatID,
	//	Text:   msg,
	//}
	reqbody2 := types.Message{Text: msg, Chat: chat,Entities: []types.MessageEntity{{Type: "mention",User: user}}}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqbody2)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1495263911:AAEjmxdCuPazMzeegGeZpm1RLBBMtFgx2oE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

// FInally, the main funtion starts our server on port 3000
func main() {
	fmt.Println("We up bois")
	port := os.Getenv("PORT")
	if port == "" {
		panic(fmt.Errorf("$PORT not set"))
	}
	err := http.ListenAndServe(port, http.HandlerFunc(telegram.HandleUpdate))
	if err != nil {
		return
	}
}