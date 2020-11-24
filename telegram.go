package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
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

func (m *SendMessage) SendMessage(c Config) {
	data, _ := json.Marshal(m)
	reqURL, _ := url.Parse("https://api.telegram.org/" + c.TelegramApiKey + "/sendMessage")
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(data)),
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to send message to telegram API")
	}
	defer res.Body.Close()
	log.Println(string(data), req)
	ioutil.ReadAll(res.Body)

}

func TelegramGetUpdates(rw http.ResponseWriter, request *http.Request, c Config) {
	var u UpdateResult

	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &u)

	resp, err := ParseRequest(u.Message.Text)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to parse search text")
	}

	if len(resp) == 0 {
		resp = "Got empty search result"
	}

	m := SendMessage{}
	m.Text = resp
	m.Chat_id = u.Chat.Id

	log.WithFields(log.Fields{
		"data":    m,
		"API Key": c.TelegramApiKey,
	}).Error("Sent message")
	m.SendMessage(c)
}
