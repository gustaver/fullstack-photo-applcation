package model

type User struct {
	Username string	`json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Error struct {
	StatusCode int
	Message string
}

type Photo struct {
	JpgBase64 string `json:"jpg"`
	Title string `json:"title"`
	Description string `json:"description"`
	Date string `json:"date"`
	User string `json:"user"`
}