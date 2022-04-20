package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	fmt.Println("create")

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(s.Command)
	fmt.Println(s.Text)
	fmt.Println("ok")
	params := &slack.Msg{Text: s.Text}
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
