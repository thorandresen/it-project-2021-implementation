package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct{
	env string
	server_addr string
	server_port int
	db_addr string
	db_port int
	db_username string
	db_password string
}

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

	return connector
}