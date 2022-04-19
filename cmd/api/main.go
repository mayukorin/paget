package main

import (
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
	}
	fmt.Println(s.Command)

}

func hello(c *gin.Context) {
	fmt.Println("hello")
}

func hello2(c *gin.Context) {
	fmt.Println("hello")
}

func main() {
	http.HandleFunc("/user-keyword-create", createUserKeyword)
	fmt.Println("main")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	/*
		router := gin.Default()
		router.GET("/hello", hello)
		router.POST("/add", createUserKeyword)
		router.Run()
	*/
}
