package main

import "fmt"

// Interface
type DatabaseRequester interface {
	getChallenge(int)
}

// Concrete1
type RealDatabaseRequester struct{}

func (rdr RealDatabaseRequester) getChallenge(pufID int) {
	fmt.Println(pufID)
}

// Concrete2
type FakeDatabaseRequester struct{}

func (fdr FakeDatabaseRequester) getChallenge(pufID int) {
	fmt.Println(pufID * 2)
}

// Strategy
type dbHandler struct {
	dbReq DatabaseRequester
}

func initDbHandler(dbr DatabaseRequester) *dbHandler {
	return &dbHandler{
		dbReq: dbr,
	}
}
func (dbh *dbHandler) setDbRequester(dr DatabaseRequester) {
	dbh.dbReq = dr
}

func main() {
	rdr := &RealDatabaseRequester{}
	// fdr := &FakeDatabaseRequester{}
	dbh := initDbHandler(rdr)
	dbh.dbReq.getChallenge(2)
}
