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
	// Only support POST requests
	if request.Method == "POST" {
		// Authenticate user and get token
		token, err := authentication.AuthenticateUser(request, model.MainDatabase)

		if err != nil {
			// If there was an error, send an error message
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send the token as a response
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(token)
		}
	} else {
		// Request was not a POST
		http.Error(writer, "API route only supports method POST", 400)
	}
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	// Only support GET requests
	if request.Method == "GET" {
		// Get photoarray from request
		photoArray, err := data.GetPhotos(request)
		if err != nil {
			// If there was an error, send an error message
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send the photo array as a response and Status OK
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(photoArray)
		}
	}

}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// Only support POST requests
	if request.Method == "POST" {
		// Upload photo from request
		err := data.UploadPhoto(request)
		if err != nil {
			// Error during photo upload
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send an OK as response since the photo has been uploaded (no error)
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Photo succesfully uploaded"))
		}
	}
}

func removeHandler(writer http.ResponseWriter, request *http.Request) {
	// Only support POST requests
	if request.Method == "POST" {
		// Remove photo from request
		err := data.RemovePhoto(request)
		if err != nil {
			// Error during photo upload
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send an OK as response since the photo has been removed (no error)
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Photo succesfully removed"))
		}
	}
}

func signupHandler(writer http.ResponseWriter, request *http.Request) {
	// Only support POST requests
	if request.Method == "POST" {
		// Sign up user from request
		err := authentication.SignupUser(request, model.MainDatabase)
		if err != nil {
			// Error during signup
			http.Error(writer, err.Message, err.StatusCode)
		} else {
			// Send OK as response since the user signed up successfully
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("User signup successful"))
		}
	}
}

// The main function that sets upp all the handle functions and calls ListenAndServe
func main() {
	// Set up database
	model.InitialiseDatabase("localhost")
	defer model.Database.Close()

	// Set up token map
	authentication.InitializeTokens()

	// Set up all routes
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/remove", removeHandler)
	http.HandleFunc("/signup", signupHandler)

	// Start server
	port := ":8080"
	log.Println("Starting server on port", port)
	http.ListenAndServe(port, nil)
}