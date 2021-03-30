package appsearch

import (
	"errors"
)

var (
	ErrInvalidParams       = errors.New("invalid params specified for Open(): accepted are (endpoint, [key])")
	ErrEngineDoesntExist   = errors.New("engine doesn't exist")
	ErrEngineAlreadyExists = errors.New("engine already exists")
)

var apiErrors = map[string]error{
	"Name is already taken":  ErrEngineAlreadyExists,
	"Could not find engine.": ErrEngineDoesntExist,
}
