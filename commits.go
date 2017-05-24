package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// var chatID int64
// var bot *tgbotapi.BotAPI

func main() {
	log.Println("Commits bot tarted")
	http.HandleFunc("/gitlab", gitlabHandler)
	http.HandleFunc("/github", githubHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}

func gitlabHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("received request: %s", req.URL)
	decoder := json.NewDecoder(req.Body)

	var api gitlab
	err := decoder.Decode(&api)
	if err != nil {
		panic(err)
	}

	var commits string
	for index, commit := range api.Commits {

		t, _ := time.Parse(time.RFC3339, commit.Timestamp)
		commitTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute())

		commits += fmt.Sprintf("\n%d. [%s](%s)\n_at %s by %s_\n",
			index+1, strings.Trim(commit.Message, "\n"),
			commit.URL, commitTime, commit.Author.Name)

		if len(commit.Added) != 0 {
			commits += fmt.Sprintf("`Added:` \n")
			for _, f := range commit.Added {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

		if len(commit.Modified) != 0 {
			commits += fmt.Sprintf("`Modified:` \n")
			for _, f := range commit.Modified {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

		if len(commit.Removed) != 0 {
			commits += fmt.Sprintf("`Removed:` \n")
			for _, f := range commit.Removed {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

	}

	text := fmt.Sprintf("*%s* pushed to *%s of %s* "+
		"[Compare changes](https://gitlab.com/%s/compare/%s...%s)\n"+
		"%s\n"+
		"Total commits: %d\n",
		api.UserName, strings.Split(api.Ref, "/")[2], api.Repository.Name,
		api.Project.PathWithNamespace, api.Before, api.CheckoutSha,
		commits,
		api.TotalCommitsCount)

	publish(text)

}

func githubHandler(w http.ResponseWriter, req *http.Request) {

	log.Printf("received request: %s", req.URL)
	decoder := json.NewDecoder(req.Body)

	var api github
	err := decoder.Decode(&api)
	if err != nil {
		panic(err)
	}

	var commits string
	for index, commit := range api.Commits {

		t, _ := time.Parse(time.RFC3339, commit.Timestamp)
		commitTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute())

		commits += fmt.Sprintf("\n%d. [%s](%s)\n_at %s by %s_\n",
			index+1, strings.Trim(commit.Message, "\n"),
			commit.URL, commitTime, commit.Author.Name)

		if len(commit.Added) != 0 {
			commits += fmt.Sprintf("`Added:` \n")
			for _, f := range commit.Added {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

		if len(commit.Modified) != 0 {
			commits += fmt.Sprintf("`Modified:` \n")
			for _, f := range commit.Modified {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

		if len(commit.Removed) != 0 {
			commits += fmt.Sprintf("`Removed:` \n")
			for _, f := range commit.Removed {
				commits += fmt.Sprintf("`— %s`\n", f)
			}
		}

	}

	text := fmt.Sprintf("*%s* pushed to *%s of %s* "+
		"[Compare changes](%s)\n"+
		"%s\n"+
		"Total commits: %d\n",
		api.Pusher.Name, strings.Split(api.Ref, "/")[2], api.Repository.Name,
		api.Compare,
		commits,
		len(api.Commits))

	publish(text)

}

// for debugging
func testHandler(w http.ResponseWriter, req *http.Request) {
	requestDump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(requestDump))
}

func publish(text string) {

	token := os.Getenv("TG_TOKEN")
	if token == "" {
		log.Fatal("token is empy")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Authorized on account %s\n\n", bot.Self.UserName)

	chatID, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}
