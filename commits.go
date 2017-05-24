package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var chatID int64
var token string
var bot *tgbotapi.BotAPI

func main() {

	checkENV()

	// Telegram
	var err error
	bot, err = tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s\n\n", bot.Self.UserName)

	log.Println("Commits bot tarted")
	http.HandleFunc("/gitlab", gitlabHandler)
	http.HandleFunc("/github", githubHandler)
	// http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}

func checkENV() {
	token = os.Getenv("TG_TOKEN")
	if token == "" {
		log.Fatal("TG_TOKEN is empy")
	}

	i, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatal("CHAT_ID is empy")
	}
	chatID = i
}

// for debugging
// func testHandler(w http.ResponseWriter, req *http.Request) {
// 	requestDump, _ := httputil.DumpRequest(req, true)
// 	fmt.Println(string(requestDump))
// }

func publish(text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}
