package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var testingEnv bool = false

var db databaseHandler

func main() {
	testingEnv := *flag.Bool("test",false,"set the env to test")
	flag.Parse()
	fmt.Println(testingEnv)
	if (testingEnv) {
		db = databaseFactory("")
	} else {
		db = databaseFactory("production")
	}    
	router := gin.Default()
    router.GET("/verify", verifyChallenge)
	router.GET("/challenge/:challenge",getChallenge)
	router.GET("/release/:id",releaseID)
	router.POST("/init/:id",facotoryInitPuf)
    router.Run("localhost:8080")
}

// Returns a challenge for a given PUF id 
func getChallenge(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("challenge"))
	challenge := db.DatabaseRequester.getChallenge(id)
	c.IndentedJSON(http.StatusOK, challenge)
}

// Verify a challange with a C,R from a given PUF ID
func verifyChallenge(c *gin.Context) {
	fmt.Println(c.PostForm("c"))
	id, _ := strconv.Atoi(c.PostForm("id"))
	challenge, _ := strconv.Atoi(c.PostForm("challenge"))
	response, _ := strconv.Atoi(c.PostForm("resposne"))
	verifyedStatus := db.DatabaseRequester.verifyChallenge(id,challenge,response)
    c.IndentedJSON(http.StatusOK, verifyedStatus)
}

// Release an given PUF id and sends a TWO step verification to a phone number.
func releaseID(c *gin.Context) {
	fmt.Println(c.PostForm("c"))
    c.IndentedJSON(http.StatusOK, "yes yes")
}

// Initiate a PUF id in database, used only for MOCK puf. 
func facotoryInitPuf(c *gin.Context) {
	fmt.Println(c.PostForm("c"))
    c.IndentedJSON(http.StatusOK, "yes yes")
}
