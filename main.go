package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func handleMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

func createNewLoginsdbFile() {
	logins := make(map[string]int)
	file, _ := json.MarshalIndent(logins, "", " ")
	_ = ioutil.WriteFile("loginsDB.json", file, 0644)
}

func handleResetCount(c *gin.Context) {
	createNewLoginsdbFile()
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func handleAddNewLogin(c *gin.Context) {
	source, s := c.GetQuery("source")
	if s {
		log.Println("request with source: " + source)
		logins := make(map[string]int)

		file, err := ioutil.ReadFile("loginsDB.json")
		if err != nil {
			log.Println("Error:", err, "Try one more time")
			createNewLoginsdbFile()
		}
		json.Unmarshal(file, &logins)

		logins[source] += 1

		file, _ = json.MarshalIndent(logins, "", " ")
		_ = ioutil.WriteFile("loginsDB.json", file, 0644)
	} else  {
		log.Println("request not consist source arg. Source not found." )
	}
	c.JSON(200, gin.H{
		"source": source,
	})
}

func handleWatchResult(c *gin.Context) {
	logins := make(map[string]int)
	file, err := ioutil.ReadFile("loginsDB.json")
	if err != nil {
		log.Println("Error:", err, "Try one more time")
	}
	json.Unmarshal(file, &logins)
	c.JSON(200, gin.H{
		"logins": logins,
	})
}

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", handleMainPage)
	router.GET("/countApp/resetCount", handleResetCount)
	router.GET("/countApp/watchResult", handleWatchResult)
	router.GET("/countApp/addNewLogin", handleAddNewLogin)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router.Run(":" + port)
}


