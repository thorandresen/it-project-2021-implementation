package main

// Interface
type DatabaseRequester interface {
	getChallenge(int) int
	verifyChallenge(int,int,int) bool
}

// Concrete1
type RealDatabaseRequester struct{}

func (rdr RealDatabaseRequester) getChallenge(pufID int) int {
	// TO be implemented :)
	return pufID
}
func (rdr RealDatabaseRequester) verifyChallenge(pufID int, challenge int, response int) bool {
	// TO be implemented :)
	return true
}

// Concrete stub, always returns (id * 2) 
type FakeDatabaseRequester struct{}

func (fdr FakeDatabaseRequester) getChallenge(pufID int) int {
	return pufID * 4
}
func (fdr FakeDatabaseRequester) verifyChallenge(pufID int, challenge int, response int) bool {
	if (challenge / 2 == response && pufID == challenge / 4) {
		return true
	}
	return false
}


// Strategy
type databaseHandler struct {
	DatabaseRequester DatabaseRequester
}

func initDbHandler(dbr DatabaseRequester) *databaseHandler {
	return &databaseHandler{
		DatabaseRequester: dbr,
	}
}

func databaseFactory(env string) databaseHandler {
	switch env {
	case "prodcution":
		return databaseHandler{
			DatabaseRequester: &RealDatabaseRequester{},
		}
	default:
		return databaseHandler{
			DatabaseRequester: &FakeDatabaseRequester{},
		} 
	}
}
