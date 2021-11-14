package main

type mitID_authtoken struct {
	tokenString string
}

//verify token with UUID
func verifyMyId(UUID string, token mitID_authtoken) bool {
	return true
}