package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
	"google.golang.org/grpc/metadata"
)

// immudb Requster struck, uses a immudb client for calls.
// a Context, and a server config from config.go file
// Config contains db password, ip addr, username and so on.
type ImmudbRequester struct{
	client client.ImmuClient
	context context.Context
	serverConfig ServerConfig
}

// Constructor for ImmudbRequester
func NewImmudbRequester(connector ServerConfig) (ir ImmudbRequester){
	ir.serverConfig = connector
	ir.client, _ = immuclient.NewImmuClient(client.DefaultOptions().WithPort(connector.db_port).WithAddress(connector.db_addr))
	ir.context = context.Background()

	lr, err := ir.client.Login(ir.context, []byte(connector.db_username), []byte(connector.db_password))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// set up an authenticated context that will be required in future operations
	md := metadata.Pairs("authorization", lr.Token)
	ir.context = metadata.NewOutgoingContext(context.Background(), md)
	return ir
}

// Get Challenge from database and return the value as int
func (immudbRequester ImmudbRequester) getChallenge(pufID int) int {
	command := "SELECT challenge_counter FROM devices WHERE pid = '" + strconv.Itoa(pufID) + "'"
	res, err := immudbRequester.client.SQLQuery(immudbRequester.context,command,nil,true)
	if err != nil {
		panic(err)
	}
	challenge, _ := strconv.Atoi(schema.RenderValue(res.Rows[0].Values[0].Value))
	return challenge
}

// Verify Challengede and returns bool based on verification. Always return 
// on erronous database lookups.
func (immudbRequester ImmudbRequester) verifyChallenge(pufID int, challenge int, response string) bool {
	requestChallenge := "SELECT response FROM puf_" + strconv.Itoa(pufID) + " WHERE challenge = " + strconv.Itoa(challenge)
	res, err := immudbRequester.client.SQLQuery(immudbRequester.context,requestChallenge,nil,true)
	if err != nil {
		return false
	}
	storedResponse := schema.RenderValue(res.Rows[0].Values[0].Value)
	storedResponse = storedResponse[1 : len(storedResponse)-1]


	if (storedResponse != "0" && storedResponse == response){
		if (immudbRequester.serverConfig.burn_puf_on_succes){
			requestOwnerAndState := "select owner,state from devices where pid='" + strconv.Itoa(pufID)+ "';"
			res, err := immudbRequester.client.SQLQuery(immudbRequester.context,requestOwnerAndState,nil,true)
			if err != nil {
				return false
			}
			owner := schema.RenderValue(res.Rows[0].Values[0].Value)
			state := schema.RenderValue(res.Rows[0].Values[1].Value)
			requestIncrement := "UPSERT INTO devices(pid,owner,challenge_counter,state) VALUES ('" + strconv.Itoa(pufID) + "','"+ owner +"'," + strconv.Itoa(challenge+1) + ",'"+ state + "');"
			requestBurnChallenge := "UPSERT INTO puf_" + strconv.Itoa(pufID) + "(challenge, response) VALUES (" + strconv.Itoa(challenge) +",'0');"
			immudbRequester.client.SQLExec(immudbRequester.context,requestBurnChallenge,nil)
			immudbRequester.client.SQLExec(immudbRequester.context,requestIncrement,nil)		
		}
		return true
	}
	return false
}

// Initiate database structure according to specifications
func (immudbRequester ImmudbRequester) commenceDatabase(){
	// SQL Commands for the database initiatization -- Create tables for users and devices and storage for keys
	var commands []string
	commands = append(commands, "CREATE TABLE IF NOT EXISTS devices(pid VARCHAR[252], owner VARCHAR[252], challenge_counter INTEGER, state VARCHAR, PRIMARY KEY pid);")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS users(uuid VARCHAR[252], first_name VARCHAR, last_name VARCHAR, phone_number INTEGER, email VARCHAR, PRIMARY KEY uuid);")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS user_keys(id INTEGER AUTO_INCREMENT, uuid VARCHAR[252], public_key VARCHAR[1024], PRIMARY KEY id);")
	
	// Connect to db
	opts := client.DefaultOptions()
	opts.Username = immudbRequester.serverConfig.db_username
	opts.Password = immudbRequester.serverConfig.db_password
	opts.Database = "defaultdb"
	opts.Address = immudbRequester.serverConfig.db_addr
	opts.Port = immudbRequester.serverConfig.db_port
	
	db := stdlib.OpenDB(opts)
	defer db.Close()
	
	//Run through commands 
	for _, command := range commands{
		_, _ = db.ExecContext(context.TODO(), command)
	}
}

// Initiate a PUF from a vendor with C,R pairs
// ATM initate random R based on seed which is ID of puf
// TODO: correctly init R from puf
func (immudbRequester ImmudbRequester) initiatePuf(id int){
	command := "CREATE TABLE IF NOT EXISTS puf_" + strconv.Itoa(id) + "(challenge INTEGER, response VARCHAR[1024], PRIMARY KEY challenge);"
	//params := map[string]interface{}{"id": 1}
	//create database table for PUF with CR pairs
	_, err := immudbRequester.client.SQLExec(immudbRequester.context,command,nil)
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < 50000; i++ {
		h := sha1.New()
		s := strconv.Itoa(id) + strconv.Itoa(i)
		h.Write([]byte(s))
		bs := h.Sum(nil)
		hash := fmt.Sprintf("%x",bs)


		command := "UPSERT INTO puf_" + strconv.Itoa(id) + "(challenge, response) VALUES (" + strconv.Itoa(i) + ",'" + hash + "')"
		if i % 1000 == 0 {
			fmt.Println(command)
		}
		_ , err := immudbRequester.client.SQLExec(immudbRequester.context,command,nil)	
		if err != nil {
			panic(err)
		}
	}
}

// Stores identity of a USER. Takes UUID and PK
func (immudbRequester ImmudbRequester) storeIdentity(uuid string, pk string) bool{
	// check if users key exist, if not store the key with the uuid
	if !immudbRequester.userKeyExits(uuid, pk) {
		storePKCommand := "INSERT INTO user_keys(uuid,public_key) VALUES (@uuid,@pk)"
		_, err := immudbRequester.client.SQLExec(immudbRequester.context,storePKCommand,map[string]interface{}{"uuid": uuid, "pk": pk})
		if err != nil {
			return false
		}
	}	
	// check if user exist if not store user with info from token. Token not impl yet, random data
	if !immudbRequester.userExist(uuid) {
		storeUserCommand := "UPSERT INTO users(uuid,email,first_name,last_name,phone_number) VALUES (@uname,@email,@first,@last,@number)"
		_, err := immudbRequester.client.SQLExec(immudbRequester.context,storeUserCommand,map[string]interface{}{"uname": uuid, "email": "Joe","first":"skirt","last":"skrski","number":23})
		if err != nil {
			return false
		}
	}
	return true
}

// Helper function, not part of framework - could be done easier with mongo or sql
// Check if a users exist in the database with the unqire UUID
func (immudbRequester ImmudbRequester) userExist(uuid string) bool {
	command := "SELECT uuid FROM users WHERE uuid = @uname;"
	res, err:= immudbRequester.client.SQLQuery(immudbRequester.context,command,map[string]interface{}{"uname": uuid},false)
	if err != nil {
		panic(err)
	} else {
	}
	i := 0
	for _, r := range res.Rows {
		for _, _ = range r.Values {
			i ++
		}
	}
	if i > 0 {
		return true
	}
	return false 
}

// Helper function, not part of framework - could be done easier with mongo or sql
// Check if a users key exist in the database with the unqire UUID
func (immudbRequester ImmudbRequester) userKeyExits(uuid string, key string) bool {
	command := "SELECT public_key FROM user_keys WHERE uuid = @uname AND public_key = @ukey;"
	res, err:= immudbRequester.client.SQLQuery(immudbRequester.context,command,map[string]interface{}{"uname": uuid, "ukey": key},false)
	if err != nil {
		panic(err)
	} else {
	}
	i := 0
	for _, r := range res.Rows {
		for _, _ = range r.Values {
			i ++
		}
	}
	if i > 0 {
		return true
	}
	return false 
}