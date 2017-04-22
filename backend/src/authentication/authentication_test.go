// Tests for the file authentication.go

package authentication

import (
	"testing"
	"model"
	"gopkg.in/mgo.v2/bson"
)

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	InitializeTokens()
}

// Create a fake request with an existing user and correct password
func TestAuthenticateUserValid(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user1", "password": "password1" }`,
		"POST", "http://localhost:8080/login", "")

	token, authError := AuthenticateUser(request, model.TestDatabase)
	if authError != nil {
		t.Error("Error authenticating user", authError.Message)
	}
	if token == nil || token.Token == "" {
		t.Error("Invalid token", token.Token)
	}
}

// Create a fake request with an existing user with the wrong password
func TestAuthenticateUserWrongPassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user1", "password": "password" }`,
		"POST", "http://localhost:8080/login", "")

	token, authError := AuthenticateUser(request, model.TestDatabase)
	if authError == nil {
		t.Error("No error when authenticating invalid user")
	}
	if token != nil {
		t.Error("Token received on invalid request", token)
	}
}

// Create a fake request with a non-existing user
func TestAuthenticateUserNoSuchUser(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user", "password": "password" }`,
		"POST", "http://localhost:8080/login", "")

	token, authError := AuthenticateUser(request, model.TestDatabase)
	if authError == nil {
		t.Error("No error when authenticating invalid user")
	}
	if token != nil {
		t.Error("Token received on invalid request", token)
	}
}

// Try to sign up a valid, non-existing user
func TestSignupUserValid(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user0", "password": "password0" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError != nil {
		t.Error("Signup error", signupError.Message)
	}

	// Check if the user is in the database
	usersCollection := model.Database.DB(model.TestDatabase).C("users")
	databaseUser := new(model.User)
	databaseError := usersCollection.Find(bson.M{"username": "user0"}).One(databaseUser)
	if databaseError != nil {
		t.Error("Database error", databaseError)
	}
}

// Try to sign up an existing user with the same password (make sure username field is unique in MongoDB)
func TestSignupUserExistingSamePassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user1", "password": "password1" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with existing user, same password")
	}
}

// Try to sign up an existing user with another password (make sure username field is unique in MongoDB)
func TestSignupUserExistingOtherPassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user1", "password": "password" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with existing user, another password")
	}
}

// Try to sign up a user with an empty username
func TestSignupUserEmptyUsername(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "", "password": "password" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty username")
	}
}


// Try to sign up a user with an empty password
func TestSignupUserEmptyPassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user", "password": "" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty password")
	}
}

// Try to sign up a user with both empty username and password
func TestSignupUserEmptyUsernameAndPassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "", "password": "" }`,
		"POST", "http://localhost:8080/signup", "")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty username and password")
	}
}

// Tests the decodeJSONToUser function, makes sure it creates a correct user object
func TestDecodeJSONToUserValid(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user", "password": "password" }`,
		"POST", "http://localhost:8080/signup", "")

	user, decodeError := decodeJSONToUser(request)
	if decodeError != nil {
		t.Error("Decode error", decodeError.Message)
	}
	if user.Username != "user" {
		t.Error("Username incorrect, found:", user.Username)
	}
	if user.Password != "password" {
		t.Error("Password incorrect, found:", user.Password)
	}
}

// Tests the decodeJSONToUser function with invalid fields
func TestDecodeJSONToUserInvalidFields(t *testing.T) {
	request := model.GenerateRequest(`{ "usern": "", "passwd": "password" }`,
		"POST", "http://localhost:8080/signup", "")

	_, decodeError := decodeJSONToUser(request)
	if decodeError == nil {
		t.Error("No decode error on invalid user model")
	}
}

// Tests the decodeJSONToUser function with empty username
func TestDecodeJSONToUserNoUsername(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "", "password": "password" }`,
		"POST", "http://localhost:8080/signup", "")

	_, decodeError := decodeJSONToUser(request)
	if decodeError == nil {
		t.Error("No decode error on no username")
	}
}

// Tests the decodeJSONToUser function with empty password
func TestDecodeJSONToUserNoPassword(t *testing.T) {
	request := model.GenerateRequest(`{ "username": "user", "password": "" }`,
		"POST", "http://localhost:8080/signup", "")

	_, decodeError := decodeJSONToUser(request)
	if decodeError == nil {
		t.Error("No decode error on no password")
	}
}
