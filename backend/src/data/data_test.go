// The file data_test.go contains tests for getting photos and upload/removal

package data

import (
	"model"
	"authentication"
	"testing"
	"time"
)

// Set up everything needed for the tests
func init() {
	model.SetupTestDatabase()
	authentication.InitializeTokens()
}

// Test GetPhotos with a valid user and token
func TestGetPhotosValid(t *testing.T) {
	token := authentication.GenerateToken(&model.User{ Username: "user1", Password: "password1" }, time.Minute)
	request := model.GenerateRequest(``, "POST", "http://localhost:8080/get", token.Token)

	photoArray, getError := GetPhotos(request, model.TestDatabase)
	if getError != nil {
		t.Error("getError not nil", getError.Message)
	}
	if len(photoArray) != 1 {
		t.Error("Expected lenght 1, got", len(photoArray))
	}
	if photoArray[0].JpgBase64 != "ABCDEFGHIJKLMNOPQRSTUVWXYZ" ||
		photoArray[0].Title != "This photo belongs to user1" ||
		photoArray[0].Description != "No description" ||
		photoArray[0].Date != "Today" ||
		photoArray[0].User != "user1" {
		t.Error("Not the expected photo, got", photoArray[0])
	}
}

// Test GetPhotos with an invalid token
func TestGetPhotosInvalid(t *testing.T) {
	token := &model.Token{ Token: "thisIsAnInvalidTOken" }
	request := model.GenerateRequest(``, "POST", "http://localhost:8080/get", token.Token)

	_, getError := GetPhotos(request, model.TestDatabase)
	if getError == nil {
		t.Error("getError nil when nil was expected")
	}
}

// Test UploadPhoto with a valid photo
func TestUploadPhotoValid(t *testing.T) {
	token := authentication.GenerateToken(&model.User{ Username: "user2", Password: "password2" }, time.Minute)
	request := model.GenerateRequest(
		`{
		"jpgbase64": "NEWPHOTO",
		 "title": "new title",
		 "description": "new description",
		 "date": "new date",
		 "user": "new user"
		 }`, "POST", "http://localhost:8080/get", token.Token)

	uploadError := UploadPhoto(request, model.TestDatabase)
	if uploadError != nil {
		t.Error("uploadError with valid data", uploadError.Message)
	}
}

// Test UploadPhoto with a invalid photo
func TestUploadPhotoInvalidData(t *testing.T) {
	token := authentication.GenerateToken(&model.User{ Username: "user2", Password: "password2" }, time.Minute)
	request := model.GenerateRequest(
		`{
		"jpgbase64": "",
		 "title": "new title",
		 "description": "new description",
		 "date": "new date",
		 "user": "new user"
		 }`, "POST", "http://localhost:8080/get", token.Token)

	uploadError := UploadPhoto(request, model.TestDatabase)
	if uploadError == nil {
		t.Error("uploadError not nil when nil expected")
	}
}
