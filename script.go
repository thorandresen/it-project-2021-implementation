package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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
func execute(pufId int, m *sync.Mutex) {
	var times []int64
	// var counter = 0
	for i := 0; i < 1000; i++ {
		start := time.Now()
		var c = getChallenge(pufId)
		var r = calcResponse(c, pufId)
		verify(c, pufId, r)
		elapsed := time.Since(start).Milliseconds()
		times = append(times, elapsed)
		// counter++
		// if counter%100 == 0 {
		// 	fmt.Println("_________________________________________________________________________________")
		// 	fmt.Println(times)
		// }
	}
	logTimes(times, m)
}

func requestOwnership(pufId int, m *sync.Mutex) (v bool) {
	signingString := strconv.Itoa(pufId)

	sk, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hasher := sha256.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPSS(rand.Reader, sk, crypto.SHA256, hasher.Sum(nil), nil)
	if err != nil {
		fmt.Println(err)
	}
	// pk := &sk.PublicKey
	// fmt.Print(pk)

	// fmt.Println(sigBytes)

	data, _ := json.Marshal(map[string]string{
		"signature": string(sigBytes),
		"uuid":      "sad",
		"pufId":     strconv.Itoa(pufId),
	})

	var times []int64
	// var counter = 0
	for i := 0; i < 1000; i++ {
		start := time.Now()

		resp, err := http.Post("https://ta.anrs.dk/transfer/request", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode == 200 || resp.StatusCode == 401 {
			elapsed := time.Since(start).Milliseconds()
			times = append(times, elapsed)
			// counter++
			// if counter%100 == 0 {
			// 	fmt.Println("_________________________________________________________________________________")
			// 	fmt.Println(times)
			// }
		}
	}
	// //Verify
	// pk := &sk.PublicKey
	// h := sha256.New()
	// h.Write([]byte("2"))
	// d := h.Sum(nil)
	// if rsa.VerifyPSS(pk, crypto.SHA256, d, sigBytes, nil) == nil {
	// 	fmt.Print("Skkkrrtttt")
	// } else {
	// 	fmt.Print("pahhhh :(")
	// }

	return
}

func logTimes(times []int64, m *sync.Mutex) {
	m.Lock()
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println(times)
	sum := int64(0)
	for i := 0; i < len(times); i++ {
		sum += times[i]
	}
	avg := sum / 1000
	log.Println("Average: ")
	log.Print(avg)
	log.Println("___________________________")

	m.Unlock()
}

// Main function
func main() {
	var m sync.Mutex

	var pufId = 1
	for i := 0; i < 12; i++ {
		// go execute(pufId)
		go requestOwnership(pufId, &m)
	}
	for true {
	}

}
