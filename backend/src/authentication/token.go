// The file token.go contains variables and functions for tokens, which are needed when trying to get and upload photos

package authentication

import (
	"model"
	"net/http"
	"math/rand"
	"strconv"
	"time"
)

// Variables needed for tokens
var tokenMap map[string]*model.User
var randomizer *rand.Rand
var deleteChan chan string

// Sets up the token map, randomizer and the channel for deleting expired tokens
func InitializeTokens() {
	tokenMap = make(map[string]*model.User)
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a delete channel and delete tokens on received on that channel
	deleteChan = make(chan string)
	go func() {
		for {
			select {
			case tokenToDelete := <- deleteChan:
				delete(tokenMap, tokenToDelete)
			}
		}
	}()
}

// Generates and returns new token for a given user with a specified validity time
func GenerateToken(user *model.User, validityTime time.Duration) *model.Token {
	// Try to generate a random token until the token is unique, since
	// the token is long the chances are extremely high that the token is unique
	duplicate := true
	newToken := ""
	for duplicate {
		newToken = strconv.Itoa(randomizer.Int())
		_, duplicate = tokenMap[newToken]
	}

	// Start a new goroutine to delete the token after a given timeout
	go func() {
		<-time.After(validityTime)
		deleteChan <- newToken
	}()

	// Put the token in the tokenMap with associated user
	tokenMap[newToken] = user

	// Return the newly generated token
	return &model.Token{ Token: newToken }
}

// Authenticates a token, making sure the token is present in the map of tokens
func AuthenticateToken(request *http.Request) (*model.User, *model.Error) {
	// Get token from Http header
	requestToken := request.Header.Get("Token")

	// Check if token exists in map, if not, send error
	if user, contains := tokenMap[requestToken]; contains {
		return user, nil
	} else {
		return nil, &model.Error{ StatusCode: 401, Message: "Invalid token" }
	}
}
