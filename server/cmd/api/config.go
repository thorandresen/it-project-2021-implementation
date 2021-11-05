package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct{
	env string
	addr string
	port int
	username string
	password string
}

func importConfig(path string) ServerConfig {
	if path == "" { path = "secret.yaml" }
	yfile, _ := ioutil.ReadFile("secret.yaml")
	data := make(map[interface{}]interface{})
	_ = yaml.Unmarshal(yfile, &data)

	var connector = ServerConfig{}
	connector.env = data["env"].(string)
	connector.addr = data["addr"].(string)
	connector.port = data["port"].(int)
	connector.username = data["username"].(string)
	connector.password = data["password"].(string)

	return connector
}