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

func GetPhotos(request *http.Request) ([]*model.Photo, *model.Error) {
	err := authentication.AuthenticateToken(request)
	if err != nil {
		return nil, err
	}

	// Decode the POST request
	decoder := json.NewDecoder(request.Body)
	requestUser := new(model.User)
	// TODO Find a way to handle faulty requests in decoding
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

func UploadPhoto(request *http.Request) {

}

func RemovePhoto(request *http.Request) {

}