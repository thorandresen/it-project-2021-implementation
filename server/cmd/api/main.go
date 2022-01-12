package main

import (
	"flag"
	// "fmt"
	// "log"
	// "net/http"
	// "strconv"
	// "time"
	// "github.com/gin-gonic/gin"
)

var testingEnv bool = false

var db databaseHandler

func main() {
	//Import Arguments
	configPath := *flag.String("path", "secret.yaml", "set the path to the config file")
	flag.Parse()

	//Import Config
	config := importConfig(configPath)

	db = databaseFactory(config)

	performBenchmark(db.DatabaseRequester)

	//db.DatabaseRequester.commenceDatabase()

	// file, err := os.OpenFile("request_logs.log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.Println("test")

	// //Setup Routers
	// router := gin.Default()
	// router.POST("/verify", verifyChallenge)
	// router.GET("/challenge/:challenge", getChallenge)
	// router.POST("/release/:id", releaseID)
	// router.POST("/init/:id", facotoryInitPuf)
	// router.POST("/create-user", createNewUser)
	// router.POST("/transfer/request", requestTransfer)
	// router.POST("/transfer/accept", acceptTransfer)
	// host := config.server_addr + ":" + strconv.Itoa(config.server_port)
	// fmt.Println(host)
	// router.Run(host)
}

// // Returns a challenge for a given PUF id
// func getChallenge(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("challenge"))
// 	challenge := db.DatabaseRequester.getChallenge(id)
// 	c.IndentedJSON(http.StatusOK, challenge)
// }

// // INPUT TO verifyChallenge
// type ChallengeJSON struct {
// 	Id        string `json:"id" binding:"required"`
// 	Challenge string `json:"challenge" binding:"required"`
// 	Response  string `json:"response" binding:"required"`
// }

// //Requiest timer for
// var requestVerifyTimer []int64

// func appendRequestTimeStampVerify(time int64) {
// 	if len(requestVerifyTimer) == 1000 {
// 		sum := int64(0)
// 		for i := 0; i < len(requestVerifyTimer); i++ {
// 			sum += requestVerifyTimer[i]
// 		}
// 		avg := sum / 1000

// 		var avgFloat float64 = float64(avg) / 1000000

// 		log.Printf("REQUEST VERIFYCHALLENGE: Average latency last 1000 request: %fms", avgFloat)
// 		requestTimes = nil
// 	} else {
// 		requestTimes = append(requestVerifyTimer, time)
// 	}
// }

// // Verify a challange with a C,R from a given PUF ID
// func verifyChallenge(c *gin.Context) {
// 	start := time.Now()

// 	data := ChallengeJSON{}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	id, _ := strconv.Atoi(data.Id)
// 	challenge, _ := strconv.Atoi(data.Challenge)
// 	verificationReponse := db.DatabaseRequester.verifyChallenge(id, challenge, data.Response)
// 	if verificationReponse {
// 		c.JSON(http.StatusOK, verificationReponse)
// 		elapsed := time.Since(start).Milliseconds()
// 		appendRequestTimeStampVerify(elapsed)
// 	} else {
// 		c.JSON(http.StatusUnauthorized, verificationReponse)
// 		elapsed := time.Since(start).Milliseconds()
// 		appendRequestTimeStampVerify(elapsed)
// 	}
// }

// // Release an given PUF id and sends a TWO step verification to a phone number.
// func releaseID(c *gin.Context) {
// 	fmt.Println(c.PostForm("c"))
// 	c.IndentedJSON(http.StatusOK, "yes yes")
// }

// // Initiate a PUF id in database, used only for MOCK puf.
// func facotoryInitPuf(c *gin.Context) {
// 	puf_id, _ := strconv.Atoi(c.Param("id"))
// 	db.DatabaseRequester.initiatePuf(puf_id)
// 	c.IndentedJSON(http.StatusOK, "yes yes")
// }

// type ConfirmBuyerStuct struct {
// 	UUID      string `json:"uuid" binding:"required"`
// 	PufID     string `json:"pufID" binding:"required"`
// 	Signature string `json:"signature" binding:"required"`
// }

// //Requiest timer for
// var requestTimes []int64

// func appendRequestTimeStamp(time int64) {
// 	if len(requestTimes) == 1000 {
// 		sum := int64(0)
// 		for i := 0; i < len(requestTimes); i++ {
// 			sum += requestTimes[i]
// 		}
// 		avg := sum / 1000

// 		var avgFloat float64 = float64(avg) / 1000000

// 		log.Printf("REQUEST TRANSFER: Average latency last 1000 request: %fms", avgFloat)
// 		requestTimes = nil
// 	} else {
// 		requestTimes = append(requestTimes, time)
// 	}
// }

// // Request an transfer of Ownsership
// func requestTransfer(c *gin.Context) {

// 	start := time.Now()

// 	data := ConfirmBuyerStuct{}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	anser := db.DatabaseRequester.confirmBuyer(data.UUID, data.Signature, data.PufID)
// 	if anser {
// 		fmt.Println("verified")
// 		c.JSON(http.StatusOK, data)

// 		elapsed := time.Since(start).Milliseconds()
// 		appendRequestTimeStamp(elapsed)
// 		//log.Printf("verified took %d", elapsed)

// 		return
// 	} else {
// 		fmt.Println("bad signature")
// 		c.JSON(http.StatusUnauthorized, data)

// 		elapsed := time.Since(start)
// 		appendRequestTimeStamp(int64(elapsed))
// 		//log.Printf("bad signature took %s", elapsed)

// 		return
// 	}

// }

// // Accept an transfer of Ownsership
// func acceptTransfer(c *gin.Context) {}

// type CreateUserStructure struct {
// 	UUID      string `json:"uuid" binding:"required"`
// 	Token     string `json:"mitIdToken" binding:"required"`
// 	PublicKey string `json:"publicKey" binding:"required"`
// }

// // Create a new user
// func createNewUser(c *gin.Context) {
// 	data := CreateUserStructure{}
// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var validationToken = mitID_authtoken{data.Token}
// 	if verifyMyId(data.UUID, validationToken) {
// 		success := db.DatabaseRequester.storeIdentity(data.UUID, data.PublicKey)
// 		if !success {
// 			c.JSON(http.StatusUnauthorized, data)
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, data)
// }
