package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/orijtech/arxiv/v1"
	"github.com/slack-go/slack"
)

var wg sync.WaitGroup // 行儀良くない？

func deliveryPaper(memberId string) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channel, _, _, err := api.OpenConversation(
		&(slack.OpenConversationParameters{
			ReturnIM: true,
			Users:    []string{memberId},
		}),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	resChan, cancel, err := arxiv.Search(context.Background(), &arxiv.Query{
		Filters: []*arxiv.Filter{
			{
				Op: arxiv.OpAnd,
				Fields: []*arxiv.Field{
					{Title: "deep learning"},
					{Title: "CSI"},
					{Title: "feedback"},
				},
			},
		},
		MaxResultsPerPage: 3,
		SortBy:            arxiv.SortBySubmittedDate,
		PageNumber:        0,
		MaxPageNumber:     1,
	})

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	for resPage := range resChan {
		if err := resPage.Err; err != nil {
			fmt.Printf("#%d err: %v", resPage.PageNumber, err)
			continue
		}
		feed := resPage.Feed
		for _, entry := range feed.Entry {
			_, _, err := api.PostMessage(
				channel.ID, // 構造体の埋め込み
				slack.MsgOptionText(entry.ID, false),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}
		if resPage.PageNumber >= 2 {
			cancel()
		}
	}

	_, _, err = api.CloseConversation(channel.ID)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	wg.Done()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

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
		wg.Add(1)
		go deliveryPaper(userSlackIdSlice[i+1])
		wg.Wait()
	}
}
