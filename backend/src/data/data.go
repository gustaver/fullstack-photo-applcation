package data

import (
	"net/http"
	"authentication"
	"model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

func GetPhotos(request *http.Request) ([]*model.Photo, *model.Error) {
	// Check that the request is properly authenticated
	err, requestUser := authentication.AuthenticateToken(request)
	if err != nil {
		return nil, err
	}

	// Photo array to be populated
	photoArray := []*model.Photo{}
	// Get the photos collection from the database
	photosCollection := model.Database.DB("main").C("photos")
	// Query database for all photos which are from the authenticated requestUser and populate Photo array with result
	photoArrayError := photosCollection.Find(bson.M{"user": requestUser.Username}).All(&photoArray)
	if photoArrayError != nil {
		// TODO: This might need to be more specific
		// Error populating photo array
		return nil, &model.Error{400, "Bad request"}
	}

	// No error occurred, return photo array (might be empty)
	return photoArray, nil
}

func UploadPhoto(request *http.Request) (*model.Error) {
	// Check that the request is properly authenticated
	authenticationError, requestUser := authentication.AuthenticateToken(request)
	if authenticationError != nil {
		return authenticationError
	}

	//Decode JSON and get photo + error
	decodeError, requestPhoto := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// Set the username of the photo to be that of the authenticated requestUser
	requestPhoto.User = requestUser.Username

	// Get the photos collection from the database
	photosCollection := model.Database.DB("main").C("photos")
	// Insert new photo into database
	insertError := photosCollection.Insert(requestPhoto)
	if insertError != nil {
		// TODO: This might be too generic
		return &model.Error{ StatusCode: 400, Message: "Photo upload failed " + insertError.Error() }
	}

	// Once we get here, err should be nil if nothing went wrong, or set to some value if something did go wrong
	return nil
}

func RemovePhoto(request *http.Request) *model.Error {
	// Check that the request is properly authenticated
	err, requestUser := authentication.AuthenticateToken(request)
	if err != nil {
		// Authentication failed
		return err
	}


	// Decode JSON to get error and requestPhoto
	decodeError, requestPhoto := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// Set the requestUser of the request photo the the requestUser that has been authenticated
	requestPhoto.User = requestUser.Username

	// Get the photos collection from the database
	photosCollection := model.Database.DB("main").C("photos")
	// Remove photo from database
	removeError := photosCollection.Remove(requestPhoto)
	if removeError != nil {
		// Photo not in database
		return &model.Error{400, "Photo does not exist"}
	}

	// Once we get here, err should be nil if nothing went wrong, or set to some value if something did go wrong
	return nil
}

func decodeJSONToPhoto(request *http.Request) (*model.Error, *model.Photo) {
	// Decode the POST request to Photo struct
	decoder := json.NewDecoder(request.Body)
	requestPhoto := new(model.Photo)
	// Faulty requests simply return empty arrays
	decodeErr := decoder.Decode(requestPhoto)
	if decodeErr != nil {
		// Something went wrong in JSON decoding
		return &model.Error{400, "Bad Request, malformatted photo"}, nil
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// Validate that photo has been built correctly from JSON
	validPhotoError := requestPhoto.Validate()
	if validPhotoError != nil {
		// JSON request didn't create complete Photo struct
		return validPhotoError, nil
	}

	// Return nil (if we get here, no error has occurred) and photo from JSON
	return nil, requestPhoto
}