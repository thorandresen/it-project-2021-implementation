package main

import (
	"context"

	"github.com/codenotary/immudb/pkg/client"

	"github.com/codenotary/immudb/pkg/stdlib"
)

func initiateDatabase(ServerConfig ServerConfig) {
	// SQL Commands for the database initiatization -- Create tables for users and devices
	var commands []string
	commands = append(commands, "CREATE TABLE IF NOT EXISTS devices(pid INTEGER, owner VARCHAR, challenge_counter INTEGER, state VARCHAR, PRIMARY KEY pid);")
	commands = append(commands, "CREATE TABLE IF NOT EXISTS users(id VARCHAR[256], pk VARCHAR, first_name VARCHAR, last_name VARCHAR, phone_number INTEGER, email VARCHAR, PRIMARY KEY id)")

	// Connect to db
	opts := client.DefaultOptions()
	opts.Username = ServerConfig.db_username
	opts.Password = ServerConfig.db_password
	opts.Database = "defaultdb"
	opts.Address = ServerConfig.db_addr
	opts.Port = ServerConfig.db_port
	
	db := stdlib.OpenDB(opts)
	defer db.Close()

	//Run through commands 
	for _, command := range commands{
		_, _ = db.ExecContext(context.TODO(), command)
	}
}