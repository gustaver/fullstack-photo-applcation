// The file data.go contains functions to get/upload/remove photos

package data

import (
	"net/http"
	"authentication"
	"model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

// Get the user's photos from the database, requires a valid token
func GetPhotos(request *http.Request, database string) ([]*model.Photo, *model.Error) {
	// Check that the request is properly authenticated
	requestUser, authError := authentication.AuthenticateToken(request)
	if authError != nil {
		return nil, authError
	}

	// Photo array to be populated with the user's photos
	photoArray := []*model.Photo{}
	// Get the photos collection from the database
	photosCollection := model.Database.DB(database).C("photos")
	// Query database for all photos that belongs to the requestUser and populate the photo array with result
	photoError := photosCollection.Find(bson.M{ "user": requestUser.Username }).All(&photoArray)
	if photoError != nil {
		// Error populating the photo array
		return nil, &model.Error{ StatusCode: 400, Message: "Error querying the database" }
	}

	// No error occurred, return photo array (empty if the user has no photos)
	return photoArray, nil
}

// Upload a photo to the database, requires a valid token
func UploadPhoto(request *http.Request, database string) (*model.Error) {
	// Check that the request is properly authenticated
	requestUser, authError := authentication.AuthenticateToken(request)
	if authError != nil {
		return authError
	}

	//Decode JSON and get photo and/or error (if one occurred)
	requestPhoto, decodeError := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// Set the username of the photo to be that of the authenticated requestUser, makes sure
	// a user cannot upload photos in another user's name even though they have a valid token
	requestPhoto.User = requestUser.Username

	// Get the photos collection from the database
	photosCollection := model.Database.DB(database).C("photos")
	// Insert new photo into database
	insertError := photosCollection.Insert(requestPhoto)
	if insertError != nil {
		// If the upload fails, return an error and the error message
		return &model.Error{ StatusCode: 400, Message: "Photo upload failed " + insertError.Error() }
	}

	// If there was no error, return nil
	return nil
}

// Removes a photo from the database belonging to a user, requires a valid token
func RemovePhoto(request *http.Request, database string) *model.Error {
	// Check that the request is properly authenticated
	requestUser, authError := authentication.AuthenticateToken(request)
	if authError != nil {
		return authError
	}


	//Decode JSON and get photo and/or error (if one occurred)
	requestPhoto, decodeError := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// Set the username of the photo to be that of the authenticated requestUser, makes sure
	// a user cannot remove photos in another user's name even though they have a valid token
	requestPhoto.User = requestUser.Username

	// Get the photos collection from the database
	photosCollection := model.Database.DB(database).C("photos")
	// Remove photo from database
	removeError := photosCollection.Remove(requestPhoto)
	if removeError != nil {
		// Photo not in database, return an error
		return &model.Error{ StatusCode: 400, Message: "Photo does not exist" }
	}

	// If there was no error, return nil
	return nil
}

// Decodes a JSON request to a Photo struct
func decodeJSONToPhoto(request *http.Request) (*model.Photo, *model.Error) {
	// Decode the POST request to Photo struct
	decoder := json.NewDecoder(request.Body)
	requestPhoto := new(model.Photo)

	// Try to decode the JSON
	decodeErr := decoder.Decode(requestPhoto)
	if decodeErr != nil {
		// Something went wrong in JSON decoding, return an error
		return nil, &model.Error{ StatusCode: 400, Message: "Bad Request, malformatted photo"}
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// Validate that photo has been built correctly from JSON
	validatePhotoError := requestPhoto.Validate()
	if validatePhotoError != nil {
		// JSON request didn't create a complete Photo struct, return an error
		return nil, validatePhotoError
	}

	// If there was no error, return the Photo struct and no error
	return requestPhoto, nil
}