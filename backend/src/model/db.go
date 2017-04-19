package model

import "gopkg.in/mgo.v2"

var Database *mgo.Session

func InitialiseDatabase(databaseSource string) {
	session, err := mgo.Dial(databaseSource)
	if err != nil {
		panic(err)
	}
	Database = session
	session.SetMode(mgo.Monotonic, true)
}
