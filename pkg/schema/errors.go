package schema

import (
	"errors"
)

var (
	// Panic in result of trying to deserialize results without { raw } values (should not happen with actual appsearch.SearchResponse)
	ErrRawValue = errors.New("inner map has value other than { raw }")
)
