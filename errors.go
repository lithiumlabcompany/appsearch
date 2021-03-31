package appsearch

import (
	"errors"
)

var (
	// ErrInvalidParams Invalid params specified for appsearch.Open
	ErrInvalidParams       = errors.New("invalid params specified for Open(): accepted are (endpoint, [key])")
	// ErrEngineDoesntExist Engine you want to create already exists
	ErrEngineDoesntExist   = errors.New("engine doesn't exist")
	// ErrEngineAlreadyExists Engine you're listing doesn't exist
	ErrEngineAlreadyExists = errors.New("engine already exists")
)

var apiErrors = map[string]error{
	"Name is already taken":  ErrEngineAlreadyExists,
	"Could not find engine.": ErrEngineDoesntExist,
}
