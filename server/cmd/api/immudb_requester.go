package main


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
