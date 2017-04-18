package data

import (
	"net/http"
	"encoding/json"
)

type User struct {
	Username string	`json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Photo struct {
	JpgBase64 string `json:"jpg"`
	Title string `json:"title"`
	Description string `json:"description"`
	Date string `json:"date"`
	User string `json:"user"`
}

type Error struct {
	StatusCode int
	Message string
}

var UsersFakeDB = []*User{
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

var PhotosFakeDB = []*Photo{
	{
		"DAHSkjhjkaHDJKASHKDJ123",
		"Day at the beach",
		"Oh the weather was sooo nice",
		"Today",
		"gustave",
	},
	{
		"ASDJODIASHDij12io312",
		"ASddasd asd asd",
		"CBA",
		"Tomrrow",
		"oskar",
	},
	{
		"ASDasdkoasdlamcxc,mnmn13",
		"Dpatchaptacha",
		"CBdsdsdds",
		"lolomrrow",
		"oskar",
	},
}

// If the requested user is in the database and the password matches - return a token,
// if there was no match - return an error
func AuthenticateUser(request *http.Request) (token *Token, error *Error) {
	// TODO Make sure the request is POST
	// Decode the POST request
	decoder := json.NewDecoder(request.Body)
	requestUser := new(User)
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
			return &Token{ "validToken" }, nil
		}
	}

	// No matching user, return an error
	return nil, &Error{ 401, "Username and password combination does not exist" }
}

func GetPhotos(request *http.Request) {

}

func UploadPhoto(request *http.Request) {

}

func RemovePhoto(request *http.Request) {

}