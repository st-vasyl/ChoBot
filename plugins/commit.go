package plugins

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetCommitMessage() (string, error) {
	var result string
	client := new(http.Client)
	req, err := http.NewRequest("GET", "http://whatthecommit.com/index.txt", nil)
	if err != nil {
		log.Printf("Failed to build request. Error: %s", err)
		return result, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get responce from google. Error: %s", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to get responce from google. Error: %s", err)
		return result, err
	}
	result = string(body)
	return result, nil
}
