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
	//Import Arguments
	configPath := *flag.String("path","secret.yaml","set the path to the config file")
	flag.Parse()

	//Import Config
	config := importConfig(configPath)

	db = databaseFactory(config)

	db.DatabaseRequester.commenceDatabase()

	//Setup Routers 
	router := gin.Default()
    router.GET("/verify", verifyChallenge)
	router.GET("/challenge/:challenge",getChallenge)
	router.GET("/release/:id",releaseID)
	router.POST("/init/:id",facotoryInitPuf)
	router.POST("/create-user",createNewUser)
	router.POST("/transfer/request",requestTransfer)
	router.POST("/transfer/accept",acceptTransfer)
	host := config.server_addr + ":" + strconv.Itoa(config.server_port)
	fmt.Println(host)
    router.Run(host)
}

// Returns a challenge for a given PUF id 
func getChallenge(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("challenge"))
	challenge := db.DatabaseRequester.getChallenge(id)
	c.IndentedJSON(http.StatusOK, challenge)
}

// Verify a challange with a C,R from a given PUF ID
func verifyChallenge(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	challenge, _ := strconv.Atoi(c.PostForm("challenge"))
	response, _ := strconv.Atoi(c.PostForm("response"))
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
	puf_id, _ := strconv.Atoi(c.Param("id"))
	db.DatabaseRequester.initiatePuf(puf_id)
    c.IndentedJSON(http.StatusOK, "yes yes")
}

// Request an transfer of Ownsership
func requestTransfer(c *gin.Context) {}

// Accept an transfer of Ownsership
func acceptTransfer(c *gin.Context) {}

// Create a new user 
func createNewUser(c *gin.Context) {
	uuid := c.PostForm("uuid")
	token := c.PostForm("mitID_token")
	var validationToken = mitID_authtoken{token}
	pk := c.PostForm("public_key")
	if verifyMyId(uuid, validationToken) {
		db.DatabaseRequester.storeIdentity(uuid,pk)
	}
	c.IndentedJSON(http.StatusOK,pk)
}