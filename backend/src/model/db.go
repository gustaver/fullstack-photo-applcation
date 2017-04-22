// The file db.go contains methods and variables for the database

package model

import (
	"gopkg.in/mgo.v2"
)

var Database *mgo.Session
const MainDatabase = "main"

// Initialises the database, dials and sets up MongoDB
func InitialiseDatabase(databaseSource string) {
	session, err := mgo.Dial(databaseSource)
	if err != nil {
		panic(err)
	}
	Database = session
	session.SetMode(mgo.Monotonic, true)
	// Make sure in Safe mode
	session.SetSafe(&mgo.Safe{})
}
