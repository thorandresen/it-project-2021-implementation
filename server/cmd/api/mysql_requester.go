package main

// User: thor@localhost
// Password: admin

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	//_ "github.com/go-sql-driver/mysql"
)

// Concrete SQL Implmentation
type MySQLRequester struct {
	db      *sql.DB
	context context.Context
}

type devices struct {
	pid               string
	owner             string
	challenge_counter int
}

<<<<<<< HEAD
func main() {
	mySqlReq := NewMySQLRequester()
	mySqlReq.commenceDatabase()
	// mySqlReq.testQuery()
	// mySqlReq.initiatePuf(4)
	// mySqlReq.storeIdentity("hej3", "skrrtpapa3")
	// mySqlReq.getChallenge(4)
	// mySqlReq.createDevice("4", "hej3", "stille")
	// mySqlReq.verifyChallenge(4, 0, "af3e133428b9e25c55bc59fe534248e6a0c0f17b")
	mySqlReq.benchmark()
}
=======
// func main() {
// 	mySqlReq := NewMySQLRequester()
// 	mySqlReq.commenceDatabase()
// 	// mySqlReq.testQuery()
// 	//mySqlReq.initiatePuf(4)
// 	//mySqlReq.storeIdentity("hej3", "skrrtpapa3")
// 	mySqlReq.getChallenge(4)
// 	//mySqlReq.createDevice("4", "hej3", "stille")
// 	mySqlReq.verifyChallenge(4, 0, "af3e133428b9e25c55bc59fe534248e6a0c0f17b")
// }
>>>>>>> 472d9907bf4df5d495f394565739b15469da49cd

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
	db, err := sql.Open("mysql", "thor:admin@tcp(127.0.0.1:3306)/defaultdb?charset=utf8&autocommit=true")
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

	command := "SELECT challenge_counter FROM devices WHERE pid LIKE " + strconv.Itoa(pufID)
	err := mySqlRequester.db.QueryRow(command).Scan(&challenge)

	if err != nil {
		panic(err)
	}

	fmt.Println(challenge)
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
		fmt.Println("VERIFIED CHALLENGE")
		return true
	}

	tx.Commit()
	fmt.Println("NOT VERIFIED CHALLENGE")
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
		res, err := db.ExecContext(mySqlRequester.context, command)

		if err != nil {
			panic(err)
		} else {
			fmt.Println(res)
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
	for i := 0; i < 30000; i++ {
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
		fmt.Println("User stored")
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
	} else {
		fmt.Println(res)
	}
	var exists int

	res.Next()
	res.Scan(&exists)
	fmt.Println(exists)

	tx.Commit()

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
