package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	p "github.com/x3rmrf/chobot/plugins"
)

var c Config

func main() {
	configfile := flag.String("configfile", "/opt/chobot/chobot.conf", "--configfile=/path/to/chobot.conf")
	flag.Parse()
	c.SetConfig(*configfile)
	GetUpdates()

}

func GetUpdates() {
	http.HandleFunc("/skype/updates", SkypeGetUpdates)
	http.HandleFunc("/telegram/updates", TelegramGetUpdates)
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
		result, _ = p.GoogleSearch(searchString, "isch", c.ConfigGoogle.GoogleSearchSafe, false)

	case "gif":
		result, _ = p.GoogleSearch(searchString, "isch", c.ConfigGoogle.GoogleSearchSafe, true)

	case "video":
		result, _ = p.GoogleSearch("site:youtube.com "+searchString, "vid", c.ConfigGoogle.GoogleSearchSafe, false)

	case "gs":
		result, _ = p.GoogleSearch(searchString, "", c.ConfigGoogle.GoogleSearchSafe, false)

	case "map":
		result = p.GoogleMaps(searchString, "640x400", "sattelite", "5")

	case "commit":
		result, _ = p.GetCommitMessage()

	case "adme":
		result, _ = p.GetAdme(searchString)

	case "Cho":
		result = p.ShowHelp()

	case "help":
		result = p.ShowHelp()

	default:
		result = "Арбалет через плечо!"
	}

	return result, nil
}
