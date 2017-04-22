package authentication

import (
	"testing"
	"model"
	"strings"
	"net/http"
	"gopkg.in/mgo.v2/bson"
)

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	InitializeTokens()
}

// Generates and returns a fake request based on the parameters of the function
func generateRequest(query, method, url string) *http.Request {
	reader := strings.NewReader(query)
	request, requestError := http.NewRequest(method, url, reader)
	if requestError != nil {
		panic(requestError)
	}
	return request
}

// Create a fake request with an existing user and correct password
func TestAuthenticateUserValid(t *testing.T) {
	request := generateRequest(`{ "username": "user1", "password": "password1" }`,
		"POST", "http://localhost:8080/login")

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
	request := generateRequest(`{ "username": "user1", "password": "password" }`,
		"POST", "http://localhost:8080/login")

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
	request := generateRequest(`{ "username": "user", "password": "password" }`,
		"POST", "http://localhost:8080/login")

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
	request := generateRequest(`{ "username": "user0", "password": "password0" }`,
		"POST", "http://localhost:8080/signup")

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
	request := generateRequest(`{ "username": "user1", "password": "password1" }`,
		"POST", "http://localhost:8080/signup")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with existing user, same password")
	}
}

// Try to sign up an existing user with another password (make sure username field is unique in MongoDB)
func TestSignupUserExistingOtherPassword(t *testing.T) {
	request := generateRequest(`{ "username": "user1", "password": "password" }`,
		"POST", "http://localhost:8080/signup")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with existing user, another password")
	}
}

// Try to sign up a user with an empty username
func TestSignupUserEmptyUsername(t *testing.T) {
	request := generateRequest(`{ "username": "", "password": "password" }`,
		"POST", "http://localhost:8080/signup")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty username")
	}
}


// Try to sign up a user with an empty password
func TestSignupUserEmptyPassword(t *testing.T) {
	request := generateRequest(`{ "username": "user", "password": "" }`,
		"POST", "http://localhost:8080/signup")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty password")
	}
}

// Try to sign up a user with both empty username and password
func TestSignupUserEmptyUsernameAndPassword(t *testing.T) {
	request := generateRequest(`{ "username": "", "password": "" }`,
		"POST", "http://localhost:8080/signup")

	// Try to sign the user up
	signupError := SignupUser(request, model.TestDatabase)
	if signupError == nil {
		t.Error("No sign up error with empty username and password")
	}
}

func TestDecodeJSONToUser(t *testing.T) {
	
}
