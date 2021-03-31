package schema

type Map = map[string]interface{}

// Schema definition as map[string]Type
// "id" field is added to schema automatically for convenience.
type Definition map[string]Type

// Schema type defines 4 types of value: text (""), date (time.RFC3339), number (0) and geolocation ("0.0,0.0")
type Type = string

const (
	TypeText        = "text"
	TypeDate        = "date"
	TypeNumber      = "number"
	TypeGeolocation = "geolocation"
)
