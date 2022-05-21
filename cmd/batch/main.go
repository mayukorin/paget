package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"github.com/mayukorin/paget"
	"github.com/orijtech/arxiv/v1"
	"github.com/slack-go/slack"
)

type Paper struct {
	ID             string
	submittedMonth string
}

type Papers []Paper

func (p Papers) Len() int {
	return len(p)
}

func (p Papers) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Papers) Less(i, j int) bool {
	return p[i].submittedMonth > p[j].submittedMonth // 逆
}

var wg sync.WaitGroup // 行儀良くない？

func deliveryPaper(slackId string) {

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channel, _, _, err := api.OpenConversation(
		&(slack.OpenConversationParameters{
			ReturnIM: true,
			Users:    []string{slackId},
		}),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q\n", err)
	}

	fmt.Println("batch")

	userId, err := paget.FindUserId(db, slackId)
	if err != nil {
		return
	}

	fmt.Println(userId)
	// rows, err := db.Query("SELECT content FROM keyword JOIN user_keyword on (keyword.id = user_keyword.keyword_id) WHERE user_keyword.slack_user_id = $1", userId)
	rows, err := paget.IndexKeywordContent(db, userId)
	if err != nil {
		return
	}

	keywordSlice := []*arxiv.Field{}

	for rows.Next() {
		var keywordContent string

		if err := rows.Scan(&keywordContent); err != nil {
			fmt.Printf("keyword content cannot get:%q\n", err)
			return
		}
		keywordSlice = append(keywordSlice, &arxiv.Field{Title: keywordContent})
	}

	resChan, cancel, err := arxiv.Search(context.Background(), &arxiv.Query{
		Filters: []*arxiv.Filter{
			{
				Op:     arxiv.OpAnd,
				Fields: keywordSlice,
			},
		},
		MaxResultsPerPage: 20,
		SortBy:            arxiv.SortByRelevance,
		PageNumber:        0,
		MaxPageNumber:     1,
	})

	if err != nil {
		fmt.Printf("%s\n", err)
	}
	var papers Papers = []Paper{} // これだけvar使ってる．

	for resPage := range resChan {
		if err := resPage.Err; err != nil {
			fmt.Printf("#%d err: %v", resPage.PageNumber, err)
			continue
		}
		feed := resPage.Feed
		for _, entry := range feed.Entry {
			fmt.Println(entry.ID)
			startIndex := strings.LastIndex(entry.ID, "/")
			fmt.Println(entry.ID[startIndex+1 : startIndex+5])
			papers = append(papers, Paper{entry.ID, entry.ID[startIndex+1:startIndex+5] + entry.ID[startIndex+6:startIndex+8]})
			/*
				_, _, err := api.PostMessage(
					channel.ID, // 構造体の埋め込み
					slack.MsgOptionText(entry.ID, false),
				)
				if err != nil {
					fmt.Printf("%s\n", err)
					return
				}
			*/
		}
		if resPage.PageNumber >= 2 {
			cancel()
		}
	}
	sort.Sort(papers)
	message := ""
	for index, paper := range papers {
		if index >= 5 {
			break
		}
		message = paper.ID + "\n"
	}

	_, _, err = api.PostMessage(
		channel.ID, // 構造体の埋め込み
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	// _, _, err = api.CloseConversation(channel.ID)
	/*
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	*/
	wg.Done()
}

func main() {
	/*
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	*/

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q\n", err)
	}

	rows, err := db.Query("SELECT slack_id FROM SLACK_USER")

	if err != nil {
		fmt.Printf("error when select keyword:%q\n", err)
		return
	}

	userSlackIdSlice := []string{}
	for rows.Next() {
		var userSlackId string
		if err := rows.Scan(&userSlackId); err != nil {
			fmt.Printf("user slack id cannot get:%q\n", err)
			return
		}
		userSlackIdSlice = append(userSlackIdSlice, userSlackId)
	}

	for i := 0; i < len(userSlackIdSlice); i += 2 {
		wg.Add(1)
		go deliveryPaper(userSlackIdSlice[i])
		if i+1 < len(userSlackIdSlice) {
			wg.Add(1)
			go deliveryPaper(userSlackIdSlice[i+1])
		}
		wg.Wait()
	}
}
