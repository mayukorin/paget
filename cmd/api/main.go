package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/mayukorin/paget"
	"github.com/slack-go/slack"
)

var db *sql.DB
var dbErr error

func indexUserKeyword(w http.ResponseWriter, r *http.Request) {

	/*
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("error opening database: %q\n", err)
		}
	*/

	fmt.Println("list")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId, err := paget.FindUserId(db, s.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		rows, err := db.Query("SELECT content FROM keyword JOIN user_keyword on (keyword.id = user_keyword.keyword_id) WHERE user_keyword.slack_user_id = $1", userId)

		if err != nil {
			fmt.Printf("error when select keyword:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/

	rows, err := paget.IndexKeywordContent(db, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseMessage := ""
	for rows.Next() {
		var keywordContent string
		if err := rows.Scan(&keywordContent); err != nil {
			fmt.Printf("keyword content cannot get:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseMessage += keywordContent + ", " // 最後いらない
	}

	params := &slack.Msg{Text: "keyword の一覧： " + responseMessage}
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("okk")
	w.Write(b)
}

func deleteUserKeyword(w http.ResponseWriter, r *http.Request) {
	/*

		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("error opening database: %q\n", err)
		}

	*/
	fmt.Println("delete")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*
		var deleteKeywordId int64
		if err := db.QueryRow("SELECT id FROM keyword WHERE content = $1", s.Text).Scan(&deleteKeywordId); err != nil {
			if err != sql.ErrNoRows {
				fmt.Printf("error when select keyword:%q\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Printf("keyword canot found")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/

	deleteKeywordId, err := paget.FindKeywordId(db, s.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId, err := paget.FindUserId(db, s.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM user_keyword WHERE keyword_id = $1 and slack_user_id = $2", deleteKeywordId, userId)
	if err != nil {
		fmt.Printf("user_keyword cannot delete: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := &slack.Msg{Text: s.Text + "を検索のキーワードから削除しました！"}
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("okk")
	w.Write(b)

}

func createUserKeyword(w http.ResponseWriter, r *http.Request) {

	/*
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("error opening database: %q\n", err)
		}
	*/
	fmt.Println("create")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId, err := paget.FindOrCreateUserId(db, s.UserID, s.ChannelID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		keywordId, err := paget.FindKeywordId(db, s.Text)
		if err != nil {
			if err != sql.ErrNoRows {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if err := db.QueryRow("INSERT INTO keyword(content) values($1) RETURNING id", s.Text).Scan(&keywordId); err != nil {
					fmt.Printf("keyword canot create:%q\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
	*/
	keywordId, err := paget.FindOrCreateKeywordId(db, s.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(keywordId)
	var userKeywordId int64
	if err := db.QueryRow("INSERT INTO user_keyword(slack_user_id, keyword_id) values($1, $2) RETURNING id", userId, keywordId).Scan(&userKeywordId); err != nil {
		fmt.Printf("keyword canot create:%q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(userKeywordId)
	params := &slack.Msg{Text: s.Text + "を検索のキーワードに追加しました！"}
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}

func main() {
	db, dbErr = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if dbErr != nil {
		fmt.Printf("error opening database: %q\n", dbErr)
	}
	fmt.Println("database open")
	fmt.Printf(db.Stats().WaitDuration.String())
	http.HandleFunc("/add", createUserKeyword)
	http.HandleFunc("/delete", deleteUserKeyword)
	http.HandleFunc("/list", indexUserKeyword)
	fmt.Println("main")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	/*
		router := gin.Default()
		router.GET("/hello", hello)
		router.POST("/add", createUserKeyword)
		router.Run()
	*/
}
