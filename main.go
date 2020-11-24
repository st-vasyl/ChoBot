package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	p "./plugins"
	log "github.com/sirupsen/logrus"
)

var c Config

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	configfile := flag.String("configfile", "/opt/chobot/chobot.conf", "--configfile=/path/to/chobot.conf")
	flag.Parse()
	c.SetConfig(*configfile)
	GetUpdates()

}

func GetUpdates() {
	http.HandleFunc("/telegram/updates", func(w http.ResponseWriter, r *http.Request) {
		TelegramGetUpdates(w, r, c)
	})
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ParseRequest(message string) (string, error) {
	if message == "" {
		return "Use '@Cho !help' for getting help", nil
	}

	var result string
	command := strings.Split(message, " ")[0][1:]
	req := strings.Split(message, " ")[1:]
	searchString := strings.Join(req, " ")

	if searchString == "" {
		return p.ShowHelp(), nil
	}

	switch command {
	case "image":
		result, _ = p.GoogleSearch(searchString, "isch", c.GoogleSafe, false)

	case "gif":
		result, _ = p.GoogleSearch(searchString, "isch", c.GoogleSafe, true)

	case "video":
		result, _ = p.GoogleSearch("site:youtube.com "+searchString, "vid", c.GoogleSafe, false)

	case "gs":
		result, _ = p.GoogleSearch(searchString, "", c.GoogleSafe, false)

	case "map":
		result = p.GoogleMaps(searchString, "640x400", "sattelite", "5")

	case "commit":
		result, _ = p.GetCommitMessage()

	case "Cho":
		result = p.ShowHelp()

	case "help":
		result = p.ShowHelp()

	default:
		result = "Арбалет через плечо!"
	}

	return result, nil
}
