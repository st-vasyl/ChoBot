package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Id         int    `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
}

type Chat struct {
	Id         int    `json:"id"`
	ChatType   string `json:"type"`
	Title      string `json:"title"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type Entities struct {
	Entitiestype string `json:"type"`
	Offset       int    `json:"offset"`
	Length       int    `json:"length"`
}

type Message struct {
	Message_id int        `json:"message_id"`
	From       User       `json:"from"`
	Chat       Chat       `json:"chat"`
	Date       int        `json:"date"`
	Text       string     `json:"text"`
	Entities   []Entities `json:"entities"`
}

type UpdateMessage struct {
	Ok     bool           `json:"ok"`
	Result []UpdateResult `json:"result"`
}

type UpdateResult struct {
	Update_id int `json:"update_id"`
	Message   `json:"message"`
}

type SendMessage struct {
	Chat_id int    `json:"chat_id"`
	Text    string `json:"text"`
}

func (m *SendMessage) SendMessage() {
	data, _ := json.Marshal(m)
	client := new(http.Client)
	postUrl := "https://api.telegram.org/" + c.TelegramApiKey + "/sendMessage"

	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "chobot/0.1")
	if err != nil {
		fmt.Println("Failed to build request ")
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Failed to POST method sendMessage ")
	}
	log.Println(data, req)
	ioutil.ReadAll(resp.Body)

}

func TelegramGetUpdates(rw http.ResponseWriter, request *http.Request) {
	var u UpdateResult

	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &u)

	resp, err := ParseRequest(u.Message.Text)
	if err != nil {
		log.Printf("Failed to parse search text. Error: %s", err)
	}

	m := SendMessage{}
	m.Text = resp
	m.Chat_id = u.Chat.Id
	m.SendMessage()
}
