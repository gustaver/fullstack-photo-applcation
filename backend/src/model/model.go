package model

// Types which implement the InputValidation interface can validate themselves after JSON decoding
type InputValidation interface {
	Validate() *Error
}

type User struct {
	Username string	`json:"username"`
	Password string `json:"password"`
}

// Validation for User type, only checks that all fields are NOT EMPTY in JSON to match User struct
func (user User) Validate() *Error {
	if user.Password == "" || user.Username == "" {
		// Either Password or Username fields have not been set
		return &Error{400, "Invalid request format"}
	}
	// JSON is valid
	return nil
}

type Token struct {
	Token string `json:"token"`
}

type Error struct {
	StatusCode int
	Message string
}

type Photo struct {
	JpgBase64 string `json:"jpgbase64"`
	Title string `json:"title"`
	Description string `json:"description"`
	Date string `json:"date"`
	User string `json:"user"`
}

// Validation for Photo type, only checks that all fields are NOT EMPTY in JSON to match Photo struct
func (photo Photo) Validate() *Error {
	if photo.User == "" || photo.Date == "" || photo.Description == "" || photo.JpgBase64 == "" || photo.Title == ""{
		// Not all of the fields have been set
		return &Error{400, "Invalid request format"}
	}
	// JSON is valid
	return nil
}