package authentication

import (
	"testing"
	"model"
	"strings"
	"net/http"
)

var database string

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	InitializeTokens()
}

func TestAuthenticateUserValid(t *testing.T) {
	// Create a fake request that should work
	reader := strings.NewReader(`{"username": "user1", "password": "password1"}`)
	request, requestError := http.NewRequest("POST", "http://localhost:8080/login", reader)
	if requestError != nil {
		panic(requestError)
	}

	token, authError := AuthenticateUser(request, model.TestDatabase)
	if authError != nil {
		t.Error("Error authenticating user", authError.Message)
	}
	if token == nil || token.Token == "" {
		t.Error("Invalid token", token.Token)
	}
}

func TestAuthenticateUserInvalid(t *testing.T) {
	// Create a fake request that shouldn't work
	reader := strings.NewReader(`{"username": "user2", "password": "password1"}`)
	request, requestError := http.NewRequest("POST", "http://localhost:8080/login", reader)
	if requestError != nil {
		panic(requestError)
	}

	token, authError := AuthenticateUser(request, model.TestDatabase)
	if authError == nil {
		t.Error("No error when authenticating invalid user")
	}
	if token != nil {
		t.Error("Token received on invalid request", token)
	}
}