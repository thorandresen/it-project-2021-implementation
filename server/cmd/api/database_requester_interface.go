package main

// Interface for the Database Requster
type DatabaseRequester interface {
	// Returns a challenge based on a PUF id
	getChallenge(int) int

	// Verify a challenge based on a PUF id, Challenge, and response
	verifyChallenge(int, int, string) bool

	// Commences a database should follow specifications ::
	commenceDatabase()

	// Initiate a PUF table in the database. Allow vendor to
	// enter C,R to a puf
	initiatePuf(int)

	// Stores the identity of a user based on a private key
	// and a UUID from mitID.
	storeIdentity(string, string) bool

	// Confirms the buyer of a product
	confirmBuyer(string, string, string) bool

	// Updates the owner of a device
	updateOwner()
}

// Strategy
type databaseHandler struct {
	DatabaseRequester DatabaseRequester
}

// Factory prodcing a Database requister based on `env` variable in config file.
// Deafult value is stubrequester to make examples work.
func databaseFactory(sc ServerConfig) databaseHandler {
	switch sc.env {
	case "production":
		return databaseHandler{
			DatabaseRequester: NewImmudbRequester(sc),
		}
	case "production-mysql":
		return databaseHandler{
			DatabaseRequester: NewMySQLRequester(sc),
		}
	default:
		return databaseHandler{
			DatabaseRequester: &StubRequester{},
		}
	}
}
