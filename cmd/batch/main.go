package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"github.com/mayukorin/paget"
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

	papers := paget.SearchArxivPapers(db, userId)

	latestMatchedPaper, err := paget.FindLatestMatchedPaper(db, slackId)

	if papers[0].ID != latestMatchedPaper {
		err = paget.UpdateUserLatestMatchedPaper(db, userId, latestMatchedPaper)
		if err != nil {

			message := papers[0].ID + "\n"

			_, _, err = api.PostMessage(
				channel.ID, // 構造体の埋め込み
				slack.MsgOptionText(message, false),
			)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			return
		}
	}

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
