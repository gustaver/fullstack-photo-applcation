// The file data_test.go contains tests for getting photos and upload/removal

package data

import (
	"model"
	"authentication"
)

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	authentication.InitializeTokens()
}


// Setup database and such
func init() {
	// Start database, make sure to run "mongod --dbpath data/test" to use testing DB
	model.InitialiseDatabase("localhost")
}
