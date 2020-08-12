package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
	"log"
	"strings"
)

type SkypeUpdateMessage struct {
	MessageType 			string     		`json:"type"`
	MessageText 			string			`json:"text"`
	MessageTime				string			`json:"timestamp"`
	MessageId				int				`json:"id"`
	MessageChannelId 		string			`json:"channelId"`
	MessageServiceUrl 		string			`json:"serviceUrl"`
	UpdateMessageFrom								`json:"from"`
	UpdateMessageConversation						`json:"conversation"`
	UpdateMessageRecipient						`json:"recipient"`
	UpdateMessageEntities					`json:"entities"`

}

type UpdateMessageFrom struct {
	MessageFromId 			string			`json:"id"`
	MessageFromName 		string			`json:"name"`
}
type UpdateMessageConversation struct {
	MessageConversationId 	string			`json:"id"`
	MessageConversationIsGroup bool 		`json:"isGroup"`
}
type UpdateMessageRecipient struct {
	MessageRecipientId 		string			`json:"id"`
	MessageRecipientName 	string			`json:"name"`
}
type UpdateMessageEntities struct {
	MessageEntitiesText string 				`json:"text"`
}


func SkypeGetUpdates(rw http.ResponseWriter, request *http.Request) {
	var s SkypeUpdateMessage
	body,_ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &s)

	if s.MessageConversationIsGroup {
		log.Printf("Message %s in group %s from %s", strings.Split(s.MessageText, " ")[2:], s.MessageConversationId, s.MessageFromName)
		s.MessageText = strings.Join(strings.Split(s.MessageText, " ")[2:], " ")
	}

	resp, err := ParseRequest(s.MessageText)
	if err != nil {
		log.Printf("Failed to parse search text. Error: %s", err)
	}

	m := SkypeSendMessage{}
	m.SendContent = resp
	m.MessageConversationId = s.MessageConversationId
	m.SkypeSendMessage()

}

type SkypeSendMessage struct {
	SkypeSendContent 	  `json:"message"`
	MessageConversationId string
}
type SkypeSendContent struct {
	SendContent string `json:"content"`
}

func (m *SkypeSendMessage) SkypeSendMessage()  {
	c.SkypeAuth()

	data,_ := json.Marshal(m)

	client := new(http.Client)
	postUrl := "https://apis.skype.com/v2/conversations/" + m.MessageConversationId + "/activities"

	authString := "Bearer " + c.SkypeConfig.AccessToken
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", authString)
	req.Header.Set("User-Agent", "ChoBot/0.1")
	if err != nil {
		fmt.Println("Failed to build request ")
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Failed to POST method sendMessage ")
	}
	ioutil.ReadAll(resp.Body)
}

func (c *Config) SkypeAuth() {
	data := "client_id=" + c.SkypeAppId + "&client_secret=" + c.SkypeSecret + "&grant_type=client_credentials&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default"
	client := new(http.Client)
	postUrl := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "ChoBot/0.1")
	if err != nil {
		fmt.Println("Failed to build request ")
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Failed to POST method sendMessage ")
	}

	result,_ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(result, &c.SkypeConfig)
}

