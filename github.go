package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func githubHandler(w http.ResponseWriter, req *http.Request) {

	log.Printf("received request: %s", req.URL)
	decoder := json.NewDecoder(req.Body)

	var api github
	err := decoder.Decode(&api)
	if err != nil {
		panic(err)
	}

	if api.Zen != nil {
		log.Println("received test hook request")
		return
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
