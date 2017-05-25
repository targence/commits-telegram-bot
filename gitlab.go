package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

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

	text := fmt.Sprintf("*%s* pushed to *%s/%s* "+
		"[Compare changes](https://gitlab.com/%s/compare/%s...%s)\n"+
		"%s\n"+
		"`Total commits: %d`\n",
		api.UserName, api.Repository.Name, strings.Split(api.Ref, "/")[2],
		api.Project.PathWithNamespace, api.Before, api.CheckoutSha,
		commits,
		api.TotalCommitsCount)

	publish(text)

}
