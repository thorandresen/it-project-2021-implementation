package main

import (
	"context"
	"log"

	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {
	client, err := immuclient.NewImmuClient(client.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password and storing a token
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	// set up an authenticated context that will be required in future operations
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	log.Println("Connection established")
}
