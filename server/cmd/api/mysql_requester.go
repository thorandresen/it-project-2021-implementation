package main

// User: thor@localhost
// Password: admin

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Concrete SQL Implmentation
type MySQLRequester struct {
	db      *sql.DB
	context context.Context
}

func main() {
	mySqlReq := NewMySQLRequester()
	mySqlReq.commenceDatabase()
	// mySqlReq.initiatePuf(1)
}

func NewMySQLRequester() (sqlRequester MySQLRequester) {
	// Capture connection properties.
	// cfg := mysql.Config{
	// 	User:   "thor@localhost",
	// 	Passwd: "admin",
	// 	Net:    "tcp",
	// 	Addr:   "127.0.0.1:3306",
	// 	DBName: "recordings",
	// }
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", "thor:admin@tcp(127.0.0.1:3306)/defaultdb")
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	sqlRequester.db = db
	sqlRequester.context = context.Background()

	return sqlRequester
}

func (mySqlRequester MySQLRequester) getChallenge(pufID int) int {
	var challenge int

	command := "SELECT challenge_counter FROM devices WHERE pid =" + strconv.Itoa(pufID)
	err := mySqlRequester.db.QueryRow(command).Scan(&challenge)
	if err != nil {
		panic(err)
	}

	return challenge
}

func (mySqlRequester MySQLRequester) verifyChallenge(pufID int, challenge int, response int) bool {

	return true
}

func (mySqlRequester MySQLRequester) commenceDatabase() {
	// SQL Commands for the database initiatization -- Create tables for users and devices and storage for keys
	var commands []string
	commands = append(commands, "CREATE TABLE IF NOT EXISTS devices(pid VARCHAR[252], owner VARCHAR[252], challenge_counter INTEGER, state VARCHAR, PRIMARY KEY pid);")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS users(uuid VARCHAR[252], first_name VARCHAR, last_name VARCHAR, phone_number INTEGER, email VARCHAR, PRIMARY KEY uuid);")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS user_keys(id INTEGER AUTO_INCREMENT, uuid VARCHAR[252], public_key VARCHAR[1024], PRIMARY KEY id);")

	db := mySqlRequester.db
	// defer db.Close()

	//Run through commands
	for _, command := range commands {
		_, _ = db.ExecContext(mySqlRequester.context, command)
	}

	fmt.Println("Commenced database")
}

func (mySqlRequester MySQLRequester) initiatePuf(id int) {
	command := "CREATE TABLE IF NOT EXISTS puf_" + strconv.Itoa(id) + "(challenge INTEGER, response INTEGER, PRIMARY KEY challenge);"
	//params := map[string]interface{}{"id": 1}
	//create database table for PUF with CR pairs
	_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)
	if err != nil {
		panic(err)
	}
	r := rand.New(rand.NewSource(int64(id)))
	for i := 0; i < 10; i++ {
		command := "UPSERT INTO puf_" + strconv.Itoa(id) + "(challenge, response) VALUES (" + strconv.Itoa(i) + "," + strconv.Itoa(r.Int()) + ")"
		_, err := mySqlRequester.db.ExecContext(mySqlRequester.context, command)

		if err != nil {
			panic(err)
		}
	}
}

func (mySqlRequester MySQLRequester) storeIdentity(uuid string, pk string) {

}

func (mySqlRequester MySQLRequester) userExist(uuid string) bool {

	return true
}

func (mySqlRequester MySQLRequester) userKeyExits(uuid string, key string, sqlRequester MySQLRequester) bool {

	return true
}
