package authentication

import (
	"model"
	"net/http"
	"math/rand"
	"strconv"
	"time"
)

var tokenMap map[string]*model.User
var randomizer *rand.Rand
var deleteChan chan string

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

func generateToken(user *model.User) *model.Token {
	newToken := strconv.Itoa(randomizer.Int())

	// Start a new goroutine to delete the token after a given timeout
	go func() {
		<-time.After(time.Minute)
		deleteChan <- newToken
	}()

	// Put the token in the tokenMap with associated user
	tokenMap[newToken] = user

	// Return the newly generated token
	return &model.Token{ Token: newToken }
}

func AuthenticateToken(request *http.Request) (*model.Error, *model.User) {
	// Get token from Http header
	requestToken := request.Header.Get("Token")

	// Check if token exists in map, if not, send error
	if user, ok := tokenMap[requestToken]; ok {
		return nil, user
	} else {
		return &model.Error{ 401, "Invalid token" }, nil
	}
}
