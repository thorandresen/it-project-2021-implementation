package main
// Concrete stub, always returns (id * 2) 
type StubRequester struct{}

func (stubRequester StubRequester) getChallenge(pufID int) int {
	return pufID * 4
}
func (stubRequester StubRequester) verifyChallenge(pufID int, challenge int, response int) bool {
	if (challenge / 2 == response && pufID == challenge / 4) {
		return true
	}
	return false
}


