package main

import (
	"context"

	"github.com/codenotary/immudb/pkg/client"

	"github.com/codenotary/immudb/pkg/stdlib"
)


func initiateDatabase(ServerConfig ServerConfig) {
	opts := client.DefaultOptions()
	opts.Username = "immudb"
	opts.Password = "immudb"
	opts.Database = "defaultdb"
	opts.Address = ServerConfig.db_addr
	opts.Port = ServerConfig.db_port
	
	db := stdlib.OpenDB(opts)
	defer db.Close()
	

	_, _ = db.ExecContext(context.TODO(), "CREATE TABLE devices(pid INTEGER, count INTEGER, owner VARCHAR, state VARCHAR, PRIMARY KEY pid)")
	_, _ = db.ExecContext(context.TODO(), "CREATE TABLE users(id VARCHAR, token VARCHAR, pk VARCHAR, name VARCHAR, number INTEGER, email VARCHAR)")
}