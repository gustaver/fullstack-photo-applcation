// The file data_test.go contains tests for getting photos and upload/removal

package data

import (
	"model"
	"authentication"
	"testing"
	"time"
	"gopkg.in/mgo.v2/bson"
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
		t.Error("getError nil when NOT nil was expected")
	}
}

// Test UploadPhoto with a valid photo
func TestUploadPhotoValid(t *testing.T) {
	// Create User for test
	testUser := &model.User{ Username: "user2", Password: "password2" }
	// Authenticate User
	token := authentication.GenerateToken(testUser, time.Minute)

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

	// Get the photos collection from the database
	photosCollection := model.Database.DB(model.TestDatabase).C("photos")
	// Query database for uploaded photo
	photo := new(model.Photo)
	photoError := photosCollection.Find(bson.M{
		"jpgbase64": "NEWPHOTO",
		"title": "new title",
		"description": "new description",
		"date": "new date",
		"user": testUser.Username}).One(photo)
	if photoError != nil {
		t.Error("Error when getting photo from database", photoError.Error())
	}

	if validateError := photo.Validate(); validateError != nil {
		t.Error("Photo was not added to database", validateError.Message)
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

// Test for removing a photo which is valid ie. is in the database for a valid user
func TestRemoveValidPhoto(t *testing.T) {
	// Create User for test
	testUser := &model.User{ Username: "user3", Password: "password3" }
	// Authenticate user
	token := authentication.GenerateToken(testUser, time.Minute)
	// Create request to remove photo that we know is in the database
	request := model.GenerateRequest(
		`{
		"jpgbase64": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"title": "This photo belongs to user3",
		"description": "No description",
		"date": "Today",
		"user": "user3"
		}`, "POST", "http://localhost:8080/get", token.Token)
	// Send request to RemovePhoto
	removeError := RemovePhoto(request, model.TestDatabase)
	if removeError != nil {
		t.Error("Remove failed despite photo being valid", removeError.Message)
	}

	photo := new(model.Photo)
	// Get the photos collection from the database
	photosCollection := model.Database.DB(model.TestDatabase).C("photos")
	// Query database for photo that we just removed, should not exist now
	photoError := photosCollection.Find(bson.M{
		"jpgbase64": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"title": "This photo belongs to user3",
		"description": "No description",
		"date": "Today",
		"user": "user3",
	}).One(photo)
	if photoError == nil {
		t.Error("Querying database for photo should have returned error")
	}
	if validateError := photo.Validate(); validateError == nil {
		t.Error("Photo validation should have returned error but didn't")
	}
}

// Test for removing a photo which is NOT valid ie. does NOT exist in the database
func TestRemoveInvalidPhoto(t *testing.T) {
	// Create User for test
	testUser := &model.User{ Username: "user4", Password: "password4" }
	// Authenticate user
	token := authentication.GenerateToken(testUser, time.Minute)
	// Create request to remove photo that we know is in NOT in the database
	request := model.GenerateRequest(
		`{
		"jpgbase64": "NOT JPEG BASE64",
		"title": "This photo belongs to NOBODY",
		"description": "No description",
		"date": "Today",
		"user": "userInvalid"
		}`, "POST", "http://localhost:8080/get", token.Token)
	// Send request to RemovePhoto
	removeError := RemovePhoto(request, model.TestDatabase)
	if removeError == nil {
		t.Error("Remove succeded despite photo being invalid")
	}
}
