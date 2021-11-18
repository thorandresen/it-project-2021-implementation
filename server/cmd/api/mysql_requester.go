package main

// Concrete SQL Implmentation
type MySQLRequester struct{
}

func NewMySQLRequester(connector ServerConfig) (sql MySQLRequester){
	
	return sql
}

func (mySqlRequester MySQLRequester) getChallenge(pufID int) int {
	
	return 0
}
func (mySqlRequester MySQLRequester) verifyChallenge(pufID int, challenge int, response int) bool {

	return true
}

func (mySqlRequester MySQLRequester) commenceDatabase(){
	
}

func (mySqlRequester MySQLRequester) initiatePuf(id int){

}

func (mySqlRequester MySQLRequester) storeIdentity(uuid string, pk string){

}

func (mySqlRequester MySQLRequester) userExist(uuid string) bool {

	return true
}

func (mySqlRequester MySQLRequester) userKeyExits(uuid string, key string, ir ImmudbRequester) bool {

	return true
}