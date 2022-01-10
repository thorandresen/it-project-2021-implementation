package main

import "strconv"

// Concrete stub, always returns (id * 2)
type StubRequester struct{}

func (stubRequester StubRequester) getChallenge(pufID int) int {
	return pufID * 4
}
func (stubRequester StubRequester) verifyChallenge(pufID int, challenge int, response string) bool {
	s, _ := strconv.Atoi(response)
	if (challenge / 2 == s && pufID == challenge / 4) {
		return true
	}
	return false
}

func (stubRequester StubRequester) commenceDatabase() {
	
}
func (stubRequester StubRequester) initiatePuf(int) {
	
}

func (stubRequester StubRequester) storeIdentity(id string, pk string ) bool{
	return true
}

func (stubRequester StubRequester) confirmBuyer(x string, y string, z string) bool{
	return true
}

func (stubRequester StubRequester) updateOwner() {

}