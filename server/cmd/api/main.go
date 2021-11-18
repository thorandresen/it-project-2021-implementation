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
	router.POST("/jsontest",testJson)
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

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// Test json
func testJson(c *gin.Context){
	json := Login{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("username: %s, pw: %s",json.User,json.Password)
	c.JSON(http.StatusAccepted,&json)
}

type ChallengeJSON struct {
	Id int `json:"id" binding:"required"`
	Challenge int `json:"challenge" binding:"required"`
	Response int `json:"response" binding:"required"`
}

// Verify a challange with a C,R from a given PUF ID
func verifyChallenge(c *gin.Context) {
	data := ChallengeJSON{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("username: %d, pw: %d",data.Id,data.Challenge)
	// _ = db.DatabaseRequester.verifyChallenge(data.Id,data.Challenge,data.Response)
	c.JSON(http.StatusOK,&data)
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