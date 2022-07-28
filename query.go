package paget

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/orijtech/arxiv/v1"
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

func FindUserId(db *sql.DB, slack_user_id string) (userId int64, err error) {

	if err = db.QueryRow("SELECT id FROM slack_user WHERE slack_id = $1", slack_user_id).Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user:%q\n", err)
		}
		fmt.Printf("slack_user canot found")
	}
	return
}

func FindLatestMatchedPaper(db *sql.DB, slack_user_id string) (latestMatchedPaper string, err error) {

	if err = db.QueryRow("SELECT latest_matched_paper FROM slack_user WHERE slack_id = $1", slack_user_id).Scan(&latestMatchedPaper); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select latest_matched_paper:%q\n", err)
		}
		fmt.Printf("slack_user canot found")
	}
	return
}

func FindOrCreateUserId(db *sql.DB, slack_user_id string, channel_id string) (userId int64, err error) {
	if userId, err = FindUserId(db, slack_user_id); err != nil {
		if err == sql.ErrNoRows {
			if err = db.QueryRow("INSERT INTO slack_user(slack_id, slack_channel_id) values ($1, $2) RETURNING id", slack_user_id, channel_id).Scan(&userId); err != nil {
				fmt.Printf("slack_user canot create:%q\n", err)
			}
		}
	}
	return
}

func FindKeywordId(db *sql.DB, keyword_content string) (keywordId int64, err error) {

	if err = db.QueryRow("SELECT id FROM keyword WHERE content = $1", keyword_content).Scan(&keywordId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error whern select keyword%q\n", err)
		} else {
			fmt.Println("keyword not found")
		}
	}
	return
}

func FindOrCreateKeywordId(db *sql.DB, keyword_content string) (keywordId int64, err error) {
	if keywordId, err = FindKeywordId(db, keyword_content); err != nil {
		if err == sql.ErrNoRows {
			if err = db.QueryRow("INSERT INTO keyword(content) values($1) RETURNING id", keyword_content).Scan(&keywordId); err != nil {
				fmt.Printf("keyword canot create:%q\n", err)
			}
		}
	}
	return

}

func UpdateUserLatestMatchedPaper(db *sql.DB, slack_user_id string, latestMatchedPaper string) (err error) {

	if _, err = db.Exec("UPDATE slack_user SET latest_matched_paper = $1 WHERE slack_id = $2", latestMatchedPaper, slack_user_id); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user:%q\n", err)
		}
		fmt.Printf("latestMatchedPaper cannnot create")
	}
	return
}

func IndexKeywordContent(db *sql.DB, user_id int64) (rows *sql.Rows, err error) {
	rows, err = db.Query("SELECT content FROM keyword JOIN user_keyword on (keyword.id = user_keyword.keyword_id) WHERE user_keyword.slack_user_id = $1", user_id)
	if err != nil {
		fmt.Printf("error when select keyword:%q\n", err)
	}
	return
}

func CreateUserKeyword(db *sql.DB, user_id int64, keyword_id int64) (userKeywordId int64, err error) {
	err = db.QueryRow("INSERT INTO user_keyword(slack_user_id, keyword_id) values($1, $2) RETURNING id", user_id, keyword_id).Scan(&userKeywordId)
	if err != nil {
		fmt.Printf("keyword canot create:%q\n", err)
	}
	return
}

func DeleteUserKeyword(db *sql.DB, userId int64, keywordId int64) (err error) {
	_, err = db.Exec("DELETE FROM user_keyword WHERE keyword_id = $1 and slack_user_id = $2", keywordId, userId)
	if err != nil {
		fmt.Printf("user_keyword cannot delete: %q\n", err)
	}
	return
}

func MakeKeywordSlice(db *sql.DB, userId int64) (keywordSlice []*arxiv.Field, err error) {
	rows, err := IndexKeywordContent(db, userId)
	if err != nil {
		return
	}

	keywordSlice = []*arxiv.Field{}
	for rows.Next() {
		var keywordContent string

		if err = rows.Scan(&keywordContent); err != nil {
			fmt.Printf("keyword content cannot get:%q\n", err)
			return
		}
		keywordSlice = append(keywordSlice, &arxiv.Field{Title: keywordContent})
	}
	return
}

func FetchArxivPapersMatchedByKeywords(keywordSlice []*arxiv.Field) (papers Papers, err error) {

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
	papers = []Paper{} // これだけvar使ってる．

	for resPage := range resChan {
		if err := resPage.Err; err != nil {
			fmt.Printf("#%d err: %v", resPage.PageNumber, err)
			continue
		}
		feed := resPage.Feed
		for _, entry := range feed.Entry {
			fmt.Println(entry.ID)
			startIndex := strings.LastIndex(entry.ID, "/")
			papers = append(papers, Paper{entry.ID, entry.ID[startIndex+1:startIndex+5] + entry.ID[startIndex+6:startIndex+8]})
		}
		if resPage.PageNumber >= 2 {
			cancel()
		}
	}
	sort.Sort(papers)

	return
}

func SearchArxivPapers(db *sql.DB, userId int64) (papers Papers) {
	keywordSlice, err := MakeKeywordSlice(db, userId)
	if err != nil {
		return
	}

	papers, err = FetchArxivPapersMatchedByKeywords(keywordSlice)
	if err != nil {
		return
	}
	return

}
