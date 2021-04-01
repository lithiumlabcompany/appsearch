package schema

import (
	"errors"
)

var (
	// Panic in result of trying to deserialize results without { raw } values (should not happen with actual appsearch.SearchResponse)
	ErrRawValue = errors.New("inner map has value other than { raw }")
	// Cannot unpack normalized map to slice
	ErrCannotUnpackSlice = errors.New("cannot Unpack map to slice. use UnpackSlice")
	// Cannot unpack to map
	ErrCannotInferFromMap = errors.New("cannot infer structure from map")
)
