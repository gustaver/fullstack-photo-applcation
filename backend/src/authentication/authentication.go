package authentication

import (
	"encoding/json"
	"net/http"
	"model"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// FIXME: Authentication should be done in HTTP header (good http protocol practice)
// If the requested user is in the database and the password matches - return a token,
// if there was no match - return an error
func AuthenticateUser(request *http.Request) (*model.Token, *model.Error) {
	decodeError, requestUser := decodeJSONToUser(request)
	if decodeError != nil {
		// Error decoding, return error and nil
		return nil, &model.Error{400, "Bad request"}
	}

	// Get collection Users from database
	usersCollection := model.Database.DB("main").C("users")
	// Check if username/password combination exists in database
	databaseUser := new(model.User)
	// Query database for user with matching username, store result in databaseUser
	databaseError := usersCollection.Find(bson.M{"username": requestUser.Username}).One(databaseUser)
	if databaseError != nil {
		// TODO: Better error response based on databaseError
		// Error querying database or user does not exist
		return nil, &model.Error{401, "Invalid login credentials"}
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(databaseUser.Password), []byte(requestUser.Password))
	if passwordError != nil {
		// Incorrect password
		return nil, &model.Error{401, "Invalid login credentials"}
	}

	// No error and user match found, return token and no error
	return &model.Token{"ValidToken"}, nil
}

// FIXME: Authentication should be done in HTTP header (good http protocol practice)
func SignupUser(request *http.Request) (*model.Error) {
	decodeError, requestUser := decodeJSONToUser(request)
	if decodeError != nil {
		// Error decoding, return error
		return decodeError
	}

	// Created hashed password
	hashedPassword, encryptionError := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	if encryptionError != nil {
		// Internal error making hashed password
		return &model.Error{500, "Internal server error creating account"}
	}
	// Set requestUser (will be put into database) password to encrypted password
	requestUser.Password = string(hashedPassword)

	// Get collection Users from database
	usersCollection := model.Database.DB("main").C("users")
	// Insert user into collection
	insertError := usersCollection.Insert(requestUser)
	if insertError != nil {
		// Error in insertion (probably due to username already taken)
		return &model.Error{400, "Bad request"}
	}

	// User successfully added to database
	return nil
}

func generateToken() *model.Token {
	// TODO real generation of token based on time and valid user
	return &model.Token{"validToken" }
}

func AuthenticateToken(request *http.Request) *model.Error {
	// TODO real authentication based on time and token generation
	return nil
}

func decodeJSONToUser(request *http.Request) (*model.Error, *model.User) {
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
		return validateUserError, nil
	}

	// Return nil (no error) and User struct decoded from JSON
	return nil, requestUser
}
