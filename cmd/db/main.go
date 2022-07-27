package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q", err)
		return
	}
	/*
		const sql6 = `
		DROP TABLE USER_KEYWORD;
		`
		if _, err := db.Exec(sql6); err != nil {
			fmt.Printf("user_keyword drop 失敗: %q", err)
			return
		}
	*/
	/*

		const sql4 = `
		DROP TABLE SLACK_USER;
		`
		if _, err := db.Exec(sql4); err != nil {
			fmt.Printf("slack_user drop 失敗: %q", err)
			return
		}

		const sql5 = `
		DROP TABLE KEYWORD;
		`
		if _, err := db.Exec(sql5); err != nil {
			fmt.Printf("keyword drop 失敗: %q", err)
			return
		}

		const sql1 = `
		CREATE TABLE SLACK_USER (
			id SERIAL PRIMARY KEY,
			slack_id VARCHAR(20) NOT NULL,
			slack_channel_id VARCHAR(20) NOT NULL
		);
		`
		if _, err := db.Exec(sql1); err != nil {
			fmt.Printf("slack_user create 失敗: %q", err)
			return
		}

		const sql2 = `
		CREATE TABLE KEYWORD (
			id SERIAL PRIMARY KEY,
			content VARCHAR(20) NOT NULL
		);
		`
		if _, err := db.Exec(sql2); err != nil {
			fmt.Printf("keyword create 失敗: %q", err)
			return
		}

		const sql3 = `
		CREATE TABLE USER_KEYWORD (
			id SERIAL PRIMARY KEY,
			slack_user_id BIGINT,
			keyword_id BIGINT,
			foreign key (slack_user_id) references SLACK_USER(id),
			foreign key (keyword_id) references KEYWORD(id)
		);
		`
		if _, err := db.Exec(sql3); err != nil {
			fmt.Printf("user keyword create 失敗: %q", err)
			return
		}
	*/

	const sql4 = `
	ALTER TABLE USER ADD COLUMN latest_matched_paper TEXT;
	`
	if _, err := db.Exec(sql4); err != nil {
		fmt.Printf("user keyword create 失敗: %q", err)
		return
	}
}
