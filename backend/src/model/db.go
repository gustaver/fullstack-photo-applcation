package model

import "gopkg.in/mgo.v2"

var Database *mgo.Session

// FIXME: Look into error handling for opening DB, and how to respond to error
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
