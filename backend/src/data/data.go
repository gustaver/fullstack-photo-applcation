package data

import (
	"net/http"
	"authentication"
	"model"
	"encoding/json"
	"log"
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
	dataBaseError := photosCollection.Insert(requestPhoto)
	if dataBaseError != nil {
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
		return err
	}

	// Decode JSON to get error and requestPhoto
	decodeError, requestPhoto := decodeJSONToPhoto(request)
	if decodeError != nil {
		return decodeError
	}

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

func decodeJSONToPhoto(request *http.Request) (*model.Error, *model.Photo) {
	// Decode the POST request to Photo struct
	decoder := json.NewDecoder(request.Body)
	requestPhoto := new(model.Photo)
	// Faulty requests simply return empty arrays
	decodeErr := decoder.Decode(requestPhoto)
	if decodeErr != nil {
		return &model.Error{
			400, "Bad Request",
		}, nil
	}
	// Close JSON decoder when function returns
	defer request.Body.Close()
	log.Print(requestPhoto)

	// Return nil (if we get here, no error has occurred) and photo from JSON
	return nil, requestPhoto
}