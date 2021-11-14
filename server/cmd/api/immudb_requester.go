package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

// Concrete1
type ImmudbRequester struct{
	client client.ImmuClient
	context context.Context
	serverConfig ServerConfig
}

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

func (immudbRequester ImmudbRequester) getChallenge(pufID int) int {
	command := "SELECT challenge_counter FROM devices WHERE pid =" + strconv.Itoa(pufID)
	res, err := immudbRequester.client.SQLQuery(immudbRequester.context,command,nil,true)
	if err != nil {
		panic(err)
	}
	challenge, _ := strconv.Atoi(schema.RenderValue(res.Rows[0].Values[0].Value))
	return challenge
}
func (immudbRequester ImmudbRequester) verifyChallenge(pufID int, challenge int, response int) bool {
	requestChallenge := "SELECT response FROM puf_" + strconv.Itoa(pufID) + " WHERE challenge = " + strconv.Itoa(challenge)
	res, _ := immudbRequester.client.SQLQuery(immudbRequester.context,requestChallenge,nil,true)
	storedResponse, _ := strconv.Atoi(schema.RenderValue(res.Rows[0].Values[0].Value))
	if (storedResponse != 0 && storedResponse == response){
		// TODO increment counter in a meaningful manner
		requestBurnChallenge := "UPSERT INTO puf_" + strconv.Itoa(pufID) + "(challenge, response) VALUES (" + strconv.Itoa(challenge) +",0)"
		immudbRequester.client.SQLExec(immudbRequester.context,requestBurnChallenge,nil)	
		return true
	}
	return false
}

func (immudbRequester ImmudbRequester) commenceDatabase(){
	initiateDatabase(immudbRequester.serverConfig)
}

func (immudbRequester ImmudbRequester) initiatePuf(id int){
	command := "CREATE TABLE IF NOT EXISTS puf_" + strconv.Itoa(id) + "(challenge INTEGER, response INTEGER, PRIMARY KEY challenge);"
	//params := map[string]interface{}{"id": 1}
	//create database table for PUF with CR pairs
	_, err := immudbRequester.client.SQLExec(immudbRequester.context,command,nil)
	if err != nil {
		panic(err)
	}
	r := rand.New(rand.NewSource(int64(id)))
	for i := 0; i < 10; i++ {
		command := "UPSERT INTO puf_" + strconv.Itoa(id) + "(challenge, response) VALUES (" + strconv.Itoa(i) + "," + strconv.Itoa(r.Int()) + ")"
		immudbRequester.client.SQLExec(immudbRequester.context,command,nil)	
	}
}

func (immudbRequester ImmudbRequester) storeIdentity(uuid string, pk string){
	if userExist(uuid) {

	} else {
		// storeUserCommand := "UPSERT INTO users(id,email,first_name,last_name,phone_number)"
	}


}

func userExist(uuid string) bool {
	// checkExistanceCommand := "SELECT "
	return true
}