package main

import (
	"net/http"
	"encoding/json"
	"log"
	"authentication"
	"data"
	"model"
)

func init() {
	log.SetPrefix("PhotoFullStack back-end: ")
}

// Handles authenticating users and sends an error message or a token as the response
func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		token, err := authentication.AuthenticateUser(request)

		if err != nil {
			// If there was an error, send an error message
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send the token as a response
			json.NewEncoder(writer).Encode(token)
		}
	} else {
		// Request was not a POST
		http.Error(writer, "API route only supports method POST", 400)
	}
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	photoArray, err := data.GetPhotos(request)
	if err != nil {
		// If there was an error, send an error message
		http.Error(writer, err.Message, err.StatusCode)
	} else {
		// Send the photo array as a response
		json.NewEncoder(writer).Encode(photoArray)
	}

}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	err := data.UploadPhoto(request)
	if err != nil {
		// Error during photo upload
		http.Error(writer, err.Message, err.StatusCode)
	} else {
		// Send an OK as response since the photo has been uploaded (no error)
		writer.WriteHeader(200)
	}
}

func removeHandler(writer http.ResponseWriter, request *http.Request) {
	err := data.RemovePhoto(request)
	if err != nil {
		// Error during photo upload
		http.Error(writer, err.Message, err.StatusCode)
	} else {
		// Send an OK as response since the photo has been uploaded (no error)
		writer.WriteHeader(200)
	}
}

// The main function that sets upp all the handle functions and calls ListenAndServe
func main() {
	model.InitialiseDatabase("localhost")
	defer model.Database.Close()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/remove", removeHandler)

	port := ":8080"
	log.Println("Starting server on port", port)
	http.ListenAndServe(port, nil)
}