package authentication

import (
	"encoding/json"
	"net/http"
	"model"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// If the requested user is in the database and the password matches - return a token,
// if there was no match - return an error
func AuthenticateUser(request *http.Request, database string) (*model.Token, *model.Error) {
	// Try to decode the requested user JSON
	requestUser, decodeError := decodeJSONToUser(request)
	if decodeError != nil {
		// Error decoding, return error and no token (nil)
		return nil, &model.Error{StatusCode: 400, Message: "Authentication failed, malformatted request" }
	}

	// Get users collection from database
	usersCollection := model.Database.DB(database).C("users")
	loginError := &model.Error{StatusCode: 401, Message: "Invalid login credentials" }

	// Query database for user with matching username, store result in databaseUser
	databaseUser := new(model.User)
	databaseError := usersCollection.Find(bson.M{"username": requestUser.Username}).One(databaseUser)
	if databaseError != nil {
		// Error querying database or user does not exist, return loginError
		return nil, loginError
	}

	// Hash password in request and compare to found user's password
	passwordError := bcrypt.CompareHashAndPassword([]byte(databaseUser.Password), []byte(requestUser.Password))
	if passwordError != nil {
		// Incorrect password, return loginError
		return nil, loginError
	}

	// No error, user match found with correct password - return token and no error
	return generateToken(databaseUser), nil
}

func SignupUser(request *http.Request, database string) (*model.Error) {
	requestUser, decodeError := decodeJSONToUser(request)
	if decodeError != nil {
		// Error decoding, return error
		return decodeError
	}

	// Created hashed password
	hashedPassword, encryptionError := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	if encryptionError != nil {
		// Internal error making hashed password
		return &model.Error{ StatusCode: 500, Message: "Internal server error creating account" }
	}
	// Set requestUser (may/may not be put into database) password to encrypted password
	requestUser.Password = string(hashedPassword)

	// Get collection Users from database
	usersCollection := model.Database.DB(database).C("users")
	// Insert user into collection
	insertError := usersCollection.Insert(requestUser)
	if insertError != nil {
		// Error in insertion (probably due to username already taken)
		return &model.Error{ StatusCode: 400, Message: "Bad request" }
	}

	// User successfully added to database
	return nil
}

func decodeJSONToUser(request *http.Request) (*model.User, *model.Error) {
	// Decode the POST request
	decoder := json.NewDecoder(request.Body)
	requestUser := new(model.User)
	err := decoder.Decode(requestUser)
	if err != nil {
		// Error during decoding
		panic(err)
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// Check that JSON decoded down to proper User struct
	validateUserError := requestUser.Validate()
	if validateUserError != nil {
		// JSON did not decode down to proper User struct
		return nil, validateUserError
	}

	// Return nil (no error) and User struct decoded from JSON
	return requestUser, nil
}
