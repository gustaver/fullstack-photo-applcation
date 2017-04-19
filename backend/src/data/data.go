package data

import (
	"net/http"
	"authentication"
	"model"
	"encoding/json"
)

// FIXME: This needs to be based on authenticated user (token perhaps)
func GetPhotos(request *http.Request) ([]*model.Photo, *model.Error) {
	// Check that the request is properly authenticated
	err := authentication.AuthenticateToken(request)
	if err != nil {
		return nil, err
	}

	// Decode the POST request
	decoder := json.NewDecoder(request.Body)
	requestUser := new(model.User)
	// Faulty requests simply return empty arrays
	decodeErr := decoder.Decode(requestUser)
	if decodeErr != nil {
		panic(err)
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()


	photoArray := []*model.Photo{}

	return photoArray, nil
}

// FIXME: This needs to be based on authenticated user (token perhaps)
func UploadPhoto(request *http.Request) (*model.Error) {
	// Check that the request is properly authenticated
	authenticationError := authentication.AuthenticateToken(request)
	if authenticationError != nil {
		return authenticationError
	}

	//Decode JSON and get photo + error
	decodeError, requestPhoto := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// Get the photos collection from the database
	photosCollection := model.Database.DB("main").C("photos")
	// Insert new photo into database
	dataBaseInsertError := photosCollection.Insert(requestPhoto)
	if dataBaseInsertError != nil {
		// TODO: This might be too generic
		return &model.Error{400, "Bad Request"}
	}

	// Once we get here, err should be nil if nothing went wrong, or set to some value if something did go wrong
	return nil
}

// FIXME: This needs to be based on authenticated user (token perhaps)
func RemovePhoto(request *http.Request) *model.Error {
	// Check that the request is properly authenticated
	err := authentication.AuthenticateToken(request)
	if err != nil {
		// Authentication failed
		return err
	}


	// Decode JSON to get error and requestPhoto
	decodeError, requestPhoto := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

	// TODO: Need to check that incoming request is properly formatted (Photo struct as JSON)
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
		return &model.Error{400, "Bad Request"}, nil
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