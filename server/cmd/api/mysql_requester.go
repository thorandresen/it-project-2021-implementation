package main

// User: thor@localhost
// Password: admin

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Concrete SQL Implmentation
type MySQLRequester struct {
	db           *sql.DB
	context      context.Context
	serverConfig ServerConfig
}

type devices struct {
	pid               string
	owner             string
	challenge_counter int
}

// func main() {
// 	mySqlReq := NewMySQLRequester()
// 	mySqlReq.commenceDatabase()
// 	// mySqlReq.testQuery()
// 	// mySqlReq.initiatePuf(6)
// 	// mySqlReq.storeIdentity("testUser1", "MjE0MjEwNzMzMjMyODc1NDU3Mjk0OTU0ODg2NzMxNjgyMjkxMzM5OTk4MDkzMTc3NDcxMzI5NzYzMDMwNjQ4MDI5NjIxNTc2Mzc3NDA5MjQxOTg2MTE5Nzk1NjcwMTQ2MTExMjkyNTU0OTk2NjAxODM2Mjg5ODg2MzEyNTQzODUzMjA4MjI3NDU2Mjg4OTc0MDg0NDE5NDM2NTUzMTQzODA4MzI3NDg5ODkwMTQzNDAxMzA4MTA4NTQ4MDEzNDIwMDM0Mjk0MTkyNzQ2MDYwMTIwNTkwMTgxNjUxMzg2NTAyMTgzNDE3MTg2NzM5NjUyMTI4MDE1MTE2NzA0NDQ0MTc0Nzc4MTIwMTE3Mjc3NTUyMTk2MTg4MDM5NDcyMTMxMzU1NTUyMTEwOTg3NzM3NDg4MDA4NjUwMTY4MzY2MTQ0MTU3MzU0MTg3NTE2MjY0ODMxNjQxMjYyMjU2MTA1NDY3MTU3NzQxMDk1NTA4MDI5OTI3MDUzMzM0NDUyODE0MjkyNDQ0NTk1MTQ1NzcxMjEyMzU2NTUwMTgwODg5MDIwMTE2MjUxODU4NDg1MjMxMzQ5NjMyNTY3MTk3NjMwODgwODc0MTAwNTQ5OTkzMDA1NzAxNDc3NTIxNTQyNzc1OTEzMzIwMTIyMDgwNzg1MTI5NTkyNjQwNDcwODEyMzQzNzU4MjE3NDA5NzQ3MTYwOTI0ODQ4NjcyNTQ4NDIyODcwNTY0MDE0NjYyODA4MDMwMjgyMTQyNDQ1NDU3Njk0MzcxMTk2NjQ1MDU2NzY4OTgxNzY3NDAxMzcxNTI5NjE5NTIxMzA0MzctNjU1Mzc=")
// 	// // mySqlReq.getChallenge(4)
// 	// mySqlReq.createDevice("6", "testUser1", "stille")
// 	// mySqlReq.verifyChallenge(4, 0, "af3e133428b9e25c55bc59fe534248e6a0c0f17b")
// 	mySqlReq.confirmBuyer("testUser1", "asasas", "6")

// 	// mySqlReq.benchmark()
// }

func NewMySQLRequester(connector ServerConfig) (sqlRequester MySQLRequester) {
	// Capture connection properties.
	// cfg := mysql.Config{
	// 	User:   "thor@localhost",
	// 	Passwd: "admin",
	// 	Net:    "tcp",
	// 	Addr:   "127.0.0.1:3306",
	// 	DBName: "recordings",
	// }
	// Get a database handle.

	sqlRequester.serverConfig = connector

	connectionString := sqlRequester.serverConfig.db_username + ":" + sqlRequester.serverConfig.db_password + "@tcp(" + sqlRequester.serverConfig.db_addr + ":" + strconv.Itoa(sqlRequester.serverConfig.db_port) + ")/defaultdb?charset=utf8&autocommit=true"

	var err error
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	//fmt.Println("Connected!")

	sqlRequester.db = db
	sqlRequester.context = context.Background()

	return sqlRequester
}

func (mySqlRequester MySQLRequester) getChallenge(pufID int) int {
	var challenge int

	command := "SELECT challenge_counter FROM devices WHERE pid LIKE " + strconv.Itoa(pufID)
	err := mySqlRequester.db.QueryRow(command).Scan(&challenge)

	if err != nil {
		panic(err)
	}
	return challenge
}

func (mySqlRequester MySQLRequester) verifyChallenge(pufID int, challenge_counter int, response string) bool {

	command := "SELECT response FROM puf_" + strconv.Itoa(pufID) + " WHERE challenge LIKE " + strconv.Itoa(challenge_counter) + ";"
	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	res, err := mySqlRequester.db.Query(command)

	if err != nil {
		panic(err)
	}

	var responseres string

	res.Next()
	res.Scan(&responseres)

	if responseres == response {
		//	fmt.Println("VERIFIED CHALLENGE")
		return true
	}

	tx.Commit()
	//fmt.Println("NOT VERIFIED CHALLENGE")
	return false
}

func (mySqlRequester MySQLRequester) commenceDatabase() {
	// SQL Commands for the database initiatization -- Create tables for users and devices and storage for keys
	var commands []string
	commands = append(commands, `CREATE TABLE IF NOT EXISTS devices(pid varchar(252) primary key, owner varchar(252), challenge_counter integer, state varchar(252))`)
	commands = append(commands, "CREATE TABLE IF NOT EXISTS users(uuid varchar(252) primary key, first_name varchar(252), last_name varchar(252), phone_number integer, email varchar(252))")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS user_keys(id integer primary key AUTO_INCREMENT, uuid varchar(252), public_key varchar(1024));")
	//commands = append(commands, "DROP TABLE puf_2")

	db := mySqlRequester.db
	// defer db.Close()
	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	//Run through commands
	for _, command := range commands {
		_, err := db.ExecContext(mySqlRequester.context, command)

		if err != nil {
			panic(err)
		} else {

		}
	}

	tx.Commit()
	fmt.Println("Commenced database")
}

func (mySqlRequester MySQLRequester) testQuery() {
	// SQL Commands for the database initiatization -- Create tables for users and devices and storage for keys

	command := "SHOW TABLES LIKE 'devices';"
	var device string

	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})
	//Run through commands
	res, err := mySqlRequester.db.Query(command)

	if err != nil {
		panic(err)
	}

	res.Scan(&device)

	tx.Commit()
	fmt.Println(device)
	fmt.Println("Tested queries")
}

func (mySqlRequester MySQLRequester) initiatePuf(id int) {
	command := "CREATE TABLE IF NOT EXISTS puf_" + strconv.Itoa(id) + "(challenge integer primary key, response varchar(1024));"
	//params := map[string]interface{}{"id": 1}
	//create database table for PUF with CR pairs
	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)

	tx.Commit()

	if err != nil {
		panic(err)
	}

	// r := rand.New(rand.NewSource(int64(id)))
	for i := 0; i < 1000; i++ {
		h := sha1.New()
		s := strconv.Itoa(id) + strconv.Itoa(i)
		h.Write([]byte(s))
		bs := h.Sum(nil)
		hash := fmt.Sprintf("%x", bs)

		// fmt.Println(strconv.Itoa(r.Int()))
		command := "INSERT INTO puf_" + strconv.Itoa(id) + "(challenge, response) VALUES (" + strconv.Itoa(i) + ", '" + hash + "');"
		//fmt.Println(command)
		_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)

		if err != nil {
			panic(err)
		}
	}
	tx.Commit()

	fmt.Println("Iniated puf")
}

func (mySqlRequester MySQLRequester) storeIdentity(uuid string, pk string) bool {
	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})
	// check if users key exist, if not store the key with the uuid
	if !mySqlRequester.userKeyExits(uuid, mySqlRequester) {
		storePKCommand := "INSERT INTO user_keys(uuid, public_key) VALUES ('" + uuid + "', '" + pk + "');"
		_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, storePKCommand)
		if err != nil {
			panic(err)
			return false
		}
		//fmt.Println("User stored")
		tx.Commit()
		return true
	}

	return false
}

func (mySqlRequester MySQLRequester) userExist(uuid string) bool {
	// command := "SELECT uuid FROM users WHERE uuid = @uname;"
	// res, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command, map[string]interface{}{"uname": uuid}, false)
	// if err != nil {
	// 	panic(err)
	// } else {
	// }
	// i := 0
	// for _, r := range res.Rows {
	// 	for _, _ = range r.Values {
	// 		i++
	// 	}
	// }
	// if i > 0 {
	// 	return true
	// }
	return false
}

func (mySqlRequester MySQLRequester) userKeyExits(uuidparam string, sqlRequester MySQLRequester) bool {
	command := "SELECT EXISTS(SELECT 1 FROM user_keys WHERE uuid LIKE '" + uuidparam + "');"

	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	// res, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)
	res, err := mySqlRequester.db.Query(command)

	if err != nil {
		panic(err)
	}
	var exists int

	res.Next()
	res.Scan(&exists)
	//fmt.Println(exists)

	tx.Commit()
	res.Close()

	if exists == 1 {
		return true
	}

	return false
}

func (mySqlRequester MySQLRequester) createDevice(pid string, owner string, state string) {
	command := "INSERT INTO devices(pid, owner, challenge_counter, state) VALUES ('" + pid + "', '" + owner + "', 0, '" + state + "');"

	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)

	if err != nil {
		panic(err)
	}

	tx.Commit()
}

func (mySqlRequester MySQLRequester) benchmark() {
	var timevar string

	command := "SELECT BENCHMARK(1000000, (SELECT response FROM puf_4 WHERE challenge LIKE 0));"

	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	res, err := mySqlRequester.db.Query(command)

	if err != nil {
		panic(err)
		return
	}

	res.Next()

	res.Scan(&timevar)

	fmt.Println(timevar)

	tx.Commit()
}

func (mySqlRequester MySQLRequester) confirmBuyer(user_id string, signature string, pufID string) bool {
	sigBytes := []byte(signature)
	//Verify

	tx, _ := mySqlRequester.db.BeginTx(mySqlRequester.context, &sql.TxOptions{})

	// Check if user exists else return false bish
	// if !NewMySQLRequester().userKeyExits(user_id, mySqlRequester) {
	// 	return false
	// }

	findPK := "SELECT public_key FROM user_keys WHERE uuid LIKE " + strconv.Quote(user_id)

	encodedStr := ""
	err := mySqlRequester.db.QueryRow(findPK).Scan(&encodedStr)

	// Exit if user don't exist

	//Decode public key
	decodedStrAsByteSlice, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		panic("malformed input")
	}

	// Slice (n,e) by -
	pk_sliced := strings.Split(string(decodedStrAsByteSlice), "-")
	n := new(big.Int)
	n, _ = n.SetString(pk_sliced[0], 10)
	e, _ := strconv.Atoi(pk_sliced[1])

	// Create rsa PK
	pk := &rsa.PublicKey{n, e}

	// Verify signature of pufID
	h := sha256.New()
	h.Write([]byte(pufID))
	d := h.Sum(nil)
	if rsa.VerifyPSS(pk, crypto.SHA256, d, sigBytes, nil) == nil {
		tx.Commit()
		return true
	} else {
		tx.Commit()
		return false
	}
}

func (mySqlRequester MySQLRequester) updateOwner() {

}

func (mySqlRequester MySQLRequester) getDatabaseType() string {
	return "mysql"
}
