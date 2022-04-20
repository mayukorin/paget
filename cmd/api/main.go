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
		r, err := db.Exec("INSERT INTO slack_user(slack_id, slack_channel_id) values ($1, $2)", s.UserID, s.ChannelID)
		if err != nil {
			fmt.Printf("slack_user canot create:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("userId cannot get:%q\n", err)
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
		r, err := db.Exec("INSERT INTO keyword(content) values($1)", s.Text)
		if err != nil {
			fmt.Printf("keyword cannot create:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		keywordId, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("keywordId cannot get:%q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	res, err := db.Exec("INSERT INTO user_keyword(slack_user_id, keyword_id) values(?, ?)", userId, keywordId)
	if err != nil {
		fmt.Printf("user_keyword cannot create:%q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userKeywordId, err := res.LastInsertId()

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
