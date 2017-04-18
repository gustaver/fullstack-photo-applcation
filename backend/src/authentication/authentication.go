package authentication

import (
	"encoding/json"
	"net/http"
	"model"
)


var UsersFakeDB = []*model.User{
	{
		"gustave",
		"12345",
	},
	{
		"oskar",
		"54321",
	},
	{
		"nobody",
		"paswd",
	},
}

// If the requested user is in the database and the password matches - return a token,
// if there was no match - return an error
func AuthenticateUser(request *http.Request) (*model.Token, *model.Error) {
	// Decode the POST request
	decoder := json.NewDecoder(request.Body)
	requestUser := new(model.User)
	// TODO Find a way to handle faulty requests in decoding
	err := decoder.Decode(requestUser)
	if err != nil {
		panic(err)
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// If there's a matching user in the database, return a token
	for _, dbUser := range UsersFakeDB {
		if requestUser.Username == dbUser.Username && requestUser.Password == dbUser.Password {
			// TODO Generate a token
			return generateToken(), nil
		}
	}

	// No matching user, return an error
	return nil, &model.Error{401, "Username and password combination does not exist" }
}

func generateToken() *model.Token {
	// TODO real generation
	return &model.Token{"validToken" }
}

func AuthenticateToken(request *http.Request) *model.Error {
	// TODO real authentication
	return nil
}
