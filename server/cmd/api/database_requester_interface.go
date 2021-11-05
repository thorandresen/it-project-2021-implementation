package main

// Interface
type DatabaseRequester interface {
	getChallenge(int) int
	verifyChallenge(int,int,int) bool
}
// Strategy
type databaseHandler struct {
	DatabaseRequester DatabaseRequester
}

func databaseFactory(env string) databaseHandler {
	switch env {
	case "prodcution":
		return databaseHandler{
			DatabaseRequester: &RealDatabaseRequester{},
		}
	default:
		return databaseHandler{
			DatabaseRequester: &StubRequester{},
		} 
	}
}
