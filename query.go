package paget

import (
	"database/sql"
	"fmt"
)

func FindUserId(db *sql.DB, slack_user_id string) (userId int64, err error) {

	if err = db.QueryRow("SELECT id FROM slack_user WHERE slack_id = $1", slack_user_id).Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user:%q\n", err)
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

func IndexKeywordContent(db *sql.DB, user_id int64) (rows *sql.Rows, err error) {
	rows, err = db.Query("SELECT content FROM keyword JOIN user_keyword on (keyword.id = user_keyword.keyword_id) WHERE user_keyword.slack_user_id = $1", user_id)
	if err != nil {
		fmt.Printf("error when select keyword:%q\n", err)
	}
	return
}
