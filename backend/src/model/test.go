// The file test.go contains methods to set up a testing database

package model

import (
	"strconv"
	"gopkg.in/mgo.v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const TestDatabase = "test"

// Clears a collection in a given database database
func clearCollection(database, collection string) {
	// Remove all users in database
	Database.DB(database).C(collection).RemoveAll(nil)
}

// Sets up the test database, clears it and populates it with fake users and photos
func SetupTestDatabase() {
	// Set up the database, make sure MongoDB is running
	InitialiseDatabase("localhost")

	// Clear the database of all current users and photos
	clearCollection(TestDatabase, "users")
	clearCollection(TestDatabase, "photos")


	// Populate the database with users and photos
	populateDatabase(TestDatabase)
}

// Populates the testing database with fake users and data
func populateDatabase(database string) {
	// Get the users and photos collections
	usersCollection := Database.DB(database).C("users")
	photosCollection := Database.DB(database).C("photos")

	// Create a number of fake users and insert them into the database, create a fake photo for every user
	for i := 1; i <= 10; i++ {
		number := strconv.Itoa(i)
		createFakeUser(usersCollection, number)
		createFakePhoto(photosCollection, number)
	}
}

// Creates a fake user based on a number and inserts it into a given collection
func createFakeUser(collection *mgo.Collection, number string) {
	// Create a hashed password based on the user's number
	hashedPassword, encryptionError := bcrypt.GenerateFromPassword([]byte("password" + number), bcrypt.DefaultCost)
	if encryptionError != nil {
		// Panic if encryption fails
		panic(encryptionError.Error())
	}

	// Insert the fake user into the database
	insertError := collection.Insert(bson.M{"username": "user" + number, "password": hashedPassword})
	if insertError != nil {
		// Panic if insertion fails
		panic(insertError.Error())
	}
}

// Creates a fake photo based on a number, connects it to a fake user and inserts it into a given collection
func createFakePhoto(collection *mgo.Collection, number string) {
	photo := &Photo{
		JpgBase64: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Title: "This photo belongs to user" + number,
		Description: "No description",
		Date: "Today",
		User: "user" + number,
	}

	// Insert the fake photo into the database
	insertError := collection.Insert(photo)
	if insertError != nil {
		// Panic if insertion fails
		panic(insertError.Error())
	}
}
