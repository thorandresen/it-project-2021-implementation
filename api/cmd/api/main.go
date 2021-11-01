package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)
func main() {
    router := gin.Default()
    router.GET("/verify", verifyChallange)
	router.GET("/challenge/:challenge",getChallenge)

    router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getChallenge(c *gin.Context) {
	fmt.Println(c.Param("challenge"))
    c.IndentedJSON(http.StatusOK, "rand challenge")
}

// getAlbums responds with the list of all albums as JSON.
func verifyChallange(c *gin.Context) {
	fmt.Println(c.PostForm("c"))
    c.IndentedJSON(http.StatusOK, "yes yes")
}