package main

import (
	"net/http"
	"data"
	"encoding/json"
	"log"
)

func init() {
	log.SetPrefix("PhotoFullStack back-end: ")
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	token, err := data.AuthenticateUser(request)

	if err != nil {
		json.NewEncoder(writer).Encode(token)
	} else {
		http.Error(writer, err.Message, err.StatusCode)
	}
}

func getHandler(writer http.ResponseWriter, request *http.Request) {

}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {

}

func removeHandler(writer http.ResponseWriter, request *http.Request) {

}

// The main function that sets upp all the handle functions and calls ListenAndServe
func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/remove", removeHandler)

	port := ":8080"
	log.Println("Starting server on port", port)
	http.ListenAndServe(port, nil)
}