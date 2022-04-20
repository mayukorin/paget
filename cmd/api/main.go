package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/slack-go/slack"
)

/*
func createUserKeyword(c *gin.Context) {
	fmt.Println("create")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(s.Command)

}
*/

func deleteUserKeyword(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q\n", err)
	}

	fmt.Println("delete")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	var userId int64
	if err := db.QueryRow("SELECT id FROM slack_user WHERE slack_id = $1", s.UserID).Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("slack_user canot found")
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q\n", err)
	}

	fmt.Println("create")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var userId int64
	fmt.Println(s.UserID)
	if err := db.QueryRow("SELECT id FROM slack_user WHERE slack_id = $1", s.UserID).Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user: %q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("slack_user not found")
		if err := db.QueryRow("INSERT INTO slack_user(slack_id, slack_channel_id) values ($1, $2) RETURNING id", s.UserID, s.ChannelID).Scan(&userId); err != nil {
			fmt.Printf("slack_user canot create:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var keywordId int64
	if err := db.QueryRow("SELECT id FROM keyword WHERE content = $1", s.Text).Scan(&keywordId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error whern select keyword%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("keyword not found")

		if err := db.QueryRow("INSERT INTO keyword(content) values($1) RETURNING id", s.Text).Scan(&keywordId); err != nil {
			fmt.Printf("keyword canot create:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

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
	fmt.Println("okk")
	w.Write(b)

}

func hello(c *gin.Context) {
	fmt.Println("hello")
}

func hello2(c *gin.Context) {
	fmt.Println("hello")
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world!")
}

func main() {
	http.HandleFunc("/add", createUserKeyword)
	http.HandleFunc("/delete", deleteUserKeyword)
	http.HandleFunc("/hello", helloWorld)
	fmt.Println("main")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	/*
		router := gin.Default()
		router.GET("/hello", hello)
		router.POST("/add", createUserKeyword)
		router.Run()
	*/
}
