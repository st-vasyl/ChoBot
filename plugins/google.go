package plugins

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func UrlEncoded(searchString, searchType, searchSafe string, isAnimated bool) (string, error) {
	var uri *url.URL
	uri, _ = url.Parse("https://encrypted.google.com")
	uri.Path += "/search"
	parameters := url.Values{}
	parameters.Add("tbm", searchType)
	parameters.Add("safe", searchSafe)
	parameters.Add("q", searchString)

	if isAnimated {
		parameters.Add("tbs", "itp:animated")
	}
	uri.RawQuery = parameters.Encode()
	log.WithFields(log.Fields{
		"plugin": "google",
		"URL":    uri.String(),
	}).Info("Encoded data")
	return uri.String(), nil
}

func UrlDecoded(str string) string {
	result, _ := url.QueryUnescape(str)
	return result
}

func GoogleSearch(searchString, searchType, searchSafe string, isAnimated bool) (string, error) {
	var result []string
	client := new(http.Client)
	useragent := "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_0 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Versio  n/4.0.5 Mobile/8A293 Safari/6531.22.7"
	encodedUri, err := UrlEncoded(searchString, searchType, searchSafe, isAnimated)
	if err != nil {
		return " ", err
	}
	r := new(regexp.Regexp)

	req, err := http.NewRequest("GET", encodedUri, nil)
	req.Header.Set("User-Agent", useragent)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "google",
			"error":  err,
		}).Error("Failed to build request")
		return " ", err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "google",
			"error":  err,
		}).Error("Failed to get responce from google")
		return " ", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "google",
			"error":  err,
		}).Error("Failed to get responce from google")
		return " ", err
	}

	switch searchType {

	case "isch":
		r, _ = regexp.Compile(`var u=\'(.*)\';var h`)
	case "vid":
		r, _ = regexp.Compile(`\?q=(https://[a-zA-z%0-9./-_]*)&amp`)
	case "":
		r, _ = regexp.Compile(`\?q=(http[a-zA-Z%:0-9./_-]*)&amp;sa`)
	}

	resArr := r.FindAllSubmatch(body, 5)
	for u := range resArr {
		result = append(result, UrlDecoded(string(resArr[u][1])))
	}
	if len(result) == 0 {
		log.WithFields(log.Fields{
			"plugin": "google",
		}).Warning("Found no result")
		return "", nil
	}

	return result[rand.Intn(5)], nil
}

func GoogleMaps(searchString, mapSize, mapType, mapZoom string) string {
	var uri *url.URL
	uri, _ = url.Parse("https://maps.googleapis.com")
	uri.Path += "/maps/api/staticmap"
	parameters := url.Values{}
	parameters.Add("size", mapSize)
	parameters.Add("markers", "size:tiny|color:0xAAAAAA%7C"+searchString)
	parameters.Add("maptype", mapType)
	parameters.Add("zoom", mapZoom)
	uri.RawQuery = parameters.Encode()
	return uri.String()
}
