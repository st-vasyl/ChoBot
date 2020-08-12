package plugins

import (
	"net/http"
	"io/ioutil"
	"log"
	"regexp"
	"math/rand"
)

func GetAdme(searchString string) (string, error) {
	var result string
	var resultArray []string
	client := new(http.Client)
	req, err := http.NewRequest("GET", "https://www.adme.ru/search/?order=date&q=" + searchString, nil)
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

	body,err  := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to get responce from google. Error: %s", err)
		return result, err
	}

	r,_ := regexp.Compile(`<a href="(.*)" class="sp-results-pic">`)
	res := r.FindAllSubmatch(body, 10)
	for i := range res {
		resultArray = append(resultArray,string(res[i][1]))
	}
	result = resultArray[rand.Intn(10)]
	log.Printf("regex result: %s", result)
	return "https://www.adme.ru" + result, nil
}
