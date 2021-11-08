package main

import (
	"context"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

// Concrete1
type ImmudbRequester struct{
	client client.ImmuClient
	context context.Context
}

func NewImmudbRequester(connector ServerConfig) (ir ImmudbRequester){
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

	vtx, err := ir.client.VerifiedSet(ir.context, []byte(`hello`), []byte(`immutable world`))
	if  err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set and verified key '%s' with value '%s' at tx %d\n", []byte(`hello`), []byte(`immutable world`), vtx.Id)

	ventry, err := ir.client.VerifiedGet(ir.context, []byte(`hello`))
	if  err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully verified key '%s' with value '%s' at tx %d\n", ventry.Key, ventry.Value, ventry.Tx)
	log.Println("Connection established")

	return ir
}

func (immudbRequester ImmudbRequester) getChallenge(pufID int) int {
	ventry, err := immudbRequester.client.VerifiedGet(immudbRequester.context, []byte(``))
	if  err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sucessfully verified key '%s' with value '%s' at tx %d\n", ventry.Key, ventry.Value, ventry.Tx)
	log.Println("Connection established")
	return pufID
}
func (immudbRequester ImmudbRequester) verifyChallenge(pufID int, challenge int, response int) bool {
	// TO be implemented :)
	return true
}