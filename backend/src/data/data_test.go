package data

import (
	"model"
	"testing"
)

type TestCase struct {
	User *model.User
	Photos *[]*model.Photo
	CorrectPhotoArray *[]*model.Photo
	CorrectErrorResponse *model.Error
}

// Setup database and such
func init() {
	// Start database, make sure to run "mongod --dbpath data/test" to use testing DB
	model.InitialiseDatabase("localhost")
}
