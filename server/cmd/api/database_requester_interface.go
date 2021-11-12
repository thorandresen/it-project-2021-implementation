package main

// Interface
type DatabaseRequester interface {
	getChallenge(int) int
	verifyChallenge(int,int,int) bool
	commenceDatabase()
	initiatePuf(int)
}
// Strategy
type databaseHandler struct {
	DatabaseRequester DatabaseRequester
}

func databaseFactory(sc ServerConfig) databaseHandler {
	switch sc.env {
	case "production":
		return databaseHandler{
			DatabaseRequester: NewImmudbRequester(sc),
		}
	default:
		return databaseHandler{
			DatabaseRequester: &StubRequester{},
		} 
	}
}
