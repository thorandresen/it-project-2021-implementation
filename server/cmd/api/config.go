package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Configuration structure for the server, holds:
// env					: Status of the env i.e Production | Server | dev etc.
// server_addr 			: Address of the API i.e. localhost
// server_port			: Port of the API i.e 4000
// db_addr				: Address to database
// db_port				: port to database
// db_username			: Database access username
// db_password			: Database access password
// burn_puf_on_succes	: burn a puf reponse on succes
type ServerConfig struct{
	env string
	burn_puf_on_succes bool
	server_addr string
	server_port int
	db_addr string
	db_port int
	db_username string
	db_password string
}

// Imorting config as yaml file from a given path
// Returns instace of ServerConfig Object.
func importConfig(path string) ServerConfig {
	if path == "" { path = "secret.yaml" }
	yfile, _ := ioutil.ReadFile("secret.yaml")
	data := make(map[interface{}]interface{})
	_ = yaml.Unmarshal(yfile, &data)

	var connector = ServerConfig{}
	connector.env = data["env"].(string)
	connector.server_addr = data["server_addr"].(string)
	connector.server_port = data["server_port"].(int)
	connector.db_addr = data["db_addr"].(string)
	connector.db_port = data["db_port"].(int)
	connector.db_username = data["db_username"].(string)
	connector.db_password = data["db_password"].(string)
	connector.burn_puf_on_succes = data["burn_puf_on_succes"].(bool)

	return connector
}