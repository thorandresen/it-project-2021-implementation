package main

// User: thor@localhost
// Password: admin

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Concrete SQL Implmentation
type MySQLRequester struct {
	db *sql.DB
}

func main() {
	mySqlReq := NewMySQLRequester()
	mySqlReq.getChallenge(5)
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

	return sqlRequester
}

func (mySqlRequester MySQLRequester) getChallenge(pufID int) int {

	return 0
}
func (mySqlRequester MySQLRequester) verifyChallenge(pufID int, challenge int, response int) bool {

	return true
}

func (mySqlRequester MySQLRequester) commenceDatabase() {

}

func (mySqlRequester MySQLRequester) initiatePuf(id int) {

}

func (mySqlRequester MySQLRequester) storeIdentity(uuid string, pk string) {

}

func (mySqlRequester MySQLRequester) userExist(uuid string) bool {

	return true
}

func (mySqlRequester MySQLRequester) userKeyExits(uuid string, key string, sqlRequester MySQLRequester) bool {

	return true
}
