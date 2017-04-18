package data

import (
	"net/http"
	"authentication"
	"model"
	"encoding/json"
)

var PhotosFakeDB = []*model.Photo{
	{
		"DAHSkjhjkaHDJKASHKDJ123",
		"Day at the beach",
		"Oh the weather was sooo nice",
		"Today",
		"gustave",
	},
	{
		"ASDJODIASHDij12io312",
		"ASddasd asd asd",
		"CBA",
		"Tomrrow",
		"oskar",
	},
	{
		"ASDasdkoasdlamcxc,mnmn13",
		"Dpatchaptacha",
		"CBdsdsdds",
		"lolomrrow",
		"oskar",
	},
}

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
	for _, photo := range PhotosFakeDB {
		if photo.User == requestUser.Username {
			photoArray = append(photoArray, photo)
		}
	}
	return photoArray, nil
}

// FIXME: This needs to be based on authenticated user (token perhaps)
func UploadPhoto(request *http.Request) (*model.Error) {
	// Check that the request is properly authenticated
	err := authentication.AuthenticateToken(request)
	if err != nil {
		return err
	}

	// Decode the POST request to Photo struct
	decoder := json.NewDecoder(request.Body)
	requestPhoto := new(model.Photo)
	// Faulty requests simply return empty arrays
	decodeErr := decoder.Decode(requestPhoto)
	if decodeErr != nil {
		panic(err)
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// Add Photo to database
	PhotosFakeDB = append(PhotosFakeDB, requestPhoto)

	// Once we get here, err should be nil if nothing went wrong, or set to some value if something did go wrong
	return err
}

// FIXME: This needs to be based on authenticated user (token perhaps)
func RemovePhoto(request *http.Request) *model.Error {
	// Check that the request is properly authenticated
	err := authentication.AuthenticateToken(request)
	if err != nil {
		return err
	}

	// Decode the POST request to Photo struct
	decoder := json.NewDecoder(request.Body)
	requestPhoto := new(model.Photo)
	// Faulty requests simply return empty arrays
	decodeErr := decoder.Decode(requestPhoto)
	if decodeErr != nil {
		panic(err)
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()

	// Find photo (if it exists in DB) and remove
	for i := len(PhotosFakeDB) - 1; i >= 0; i-- {
		 currentPhoto := PhotosFakeDB[i]
		// Condition to decide if current element has to be deleted:
		if currentPhoto.Title == requestPhoto.Title && currentPhoto.User == requestPhoto.User {
			PhotosFakeDB = append(PhotosFakeDB[:i], PhotosFakeDB[i+1:]...)
		}
	}

	// Once we get here, err should be nil if nothing went wrong, or set to some value if something did go wrong
	return err
}