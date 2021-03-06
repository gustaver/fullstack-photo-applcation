// Tests for the file token.go

package authentication

import (
	"testing"
	"model"
	"time"
)

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	InitializeTokens()
}

// Make sure that InitializeTokens initializes the variables needed for tokens
func TestInitializeTokens(t *testing.T) {
	if tokenMap == nil {
		t.Error("tokenMap not initialized")
	}
	if randomizer == nil {
		t.Error("randomizer not initialized")
	}
	if deleteChan == nil {
		t.Error("deleteChan not initialized")
	}
}

// Make sure that generateToken generates a token and removes it after the timeout has expired
func TestGenerateToken(t *testing.T) {
	token := GenerateToken(&model.User{ Username: "user", Password: "password" }, time.Millisecond*500)
	if _, contains := tokenMap[token.Token]; !contains {
		t.Error("tokenMap does not contain the newly generated token")
	}
	time.Sleep(time.Second)
	if _, contains := tokenMap[token.Token]; contains {
		t.Error("tokenMap still contains the token after timeout expired", token.Token)
	}
}

// Test to authenticate token with a valid user present in the tokenMap
func TestAuthenticateTokenValid(t *testing.T) {
	token := GenerateToken(&model.User{ Username: "user1", Password: "password1" }, time.Minute)

	request := model.GenerateRequest(``,
		"POST", "http://localhost:8080/login", token.Token)

	user, tokenError := AuthenticateToken(request)
	if tokenError != nil {
		t.Error("Authentication of token returned an error", tokenError.Message)
	}
	if user.Username != "user1" {
		t.Error("Username does not match requested user", user.Username)
	}
	if user.Password != "password1" {
		t.Error("Password does not match requested user", user.Password)
	}
}

// Test to authenticate token with an invalid user (not present in the tokenMap)
func TestAuthenticateTokenInvalid(t *testing.T) {
	request := model.GenerateRequest(``,
		"POST", "http://localhost:8080/login", "thisIsAnInvalidToken")

	_, tokenError := AuthenticateToken(request)
	if tokenError == nil {
		t.Error("No tokenError returned with invalid user")
	}
}
