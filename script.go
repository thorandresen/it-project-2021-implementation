package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Importing packages
func getChallenge(pufId int) (c int) {
	resp, err := http.Get("https://ta.anrs.dk/challenge/2")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	c, err = strconv.Atoi(string(body))
	if err != nil {
		log.Fatalln(err)
	}
	return
}
func calcResponse(c int, pufId int) (r string) {
	hasher := sha1.New()
	hasher.Write([]byte(strconv.Itoa(pufId) + strconv.Itoa(c)))
	hash := hasher.Sum(nil)
	r = fmt.Sprintf("%x", hash)
	return
}
func verify(c int, pufId int, r string) (v bool) {

	data, _ := json.Marshal(map[string]string{
		"id":        strconv.Itoa(pufId),
		"challenge": strconv.Itoa(c),
		"response":  r,
	})

	resp, err := http.Post("https://ta.anrs.dk/verify", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	v, err = strconv.ParseBool(string(body))
	if err != nil {
		log.Fatalln(err)
	}
	return
}
func execute(pufId int) {
	var times []time.Duration
	var counter = 0
	for i := 0; i < 1000; i++ {
		start := time.Now()
		var c = getChallenge(pufId)
		var r = calcResponse(c, pufId)
		verify(c, pufId, r)
		elapsed := time.Since(start)
		times = append(times, elapsed)
		counter++
		if counter%100 == 0 {
			fmt.Println("_________________________________________________________________________________")
			fmt.Println(pufId)
			fmt.Println(times)
		}
	}
}

// Main function
func main() {
	var pufId = 2
	for i := 0; i < 4; i++ {
		go execute(pufId)
	}
	for true {
	}

}
