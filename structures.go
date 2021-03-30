package appsearch

import (
	"fmt"
	"strings"
)

type m = map[string]interface{}
type stringMap = map[string]string

type Page struct {
	Page int `json:"current,omitempty"`
	Size int `json:"size,omitempty"`
}

// Response for Patch or Update operations
type UpdateResponse struct {
	// Updated document ID
	ID string `json:"id"`
	// List of errors
	Errors []string `json:"errors"`
}

// Response for Patch or Update operations
type DeleteResponse struct {
	// Deleted document ID
	ID string `json:"id"`
	// Was document deleted successfully
	Deleted bool `json:"deleted"`
	// List of errors
	Errors []string `json:"errors"`
}

// Schema type defines 4 types of value: text (""), date (time.RFC3339), number (0) and geolocation ("0.0,0.0")
type SchemaType = string

const (
	SchemaTypeText        = "text"
	SchemaTypeDate        = "date"
	SchemaTypeNumber      = "number"
	SchemaTypeGeolocation = "geolocation"
)

// Schema definition as map[string]SchemaType
// "id" field of "text" type is added to schema automatically (non-standard behaviour).
type SchemaDefinition map[string]SchemaType

// Pagination metadata included in paged responses
type PaginationMeta struct {
	PageSize     int `json:"size"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current"`
	TotalResults int `json:"total_results"`
}

// Response metadata included in some responses
type ResponseMeta struct {
	Page      PaginationMeta `json:"page"`
	Alerts    []string       `json:"alerts"`
	Warnings  []string       `json:"warnings"`
	RequestID string         `json:"request_id"`
}

// ListEngines response
type EngineResponse struct {
	Meta    ResponseMeta        `json:"meta"`
	Results []EngineDescription `json:"results"`
}

// Request for CreateEngine
type CreateEngineRequest struct {
	Name     string `json:"name"`
	Language string `json:"language,omitempty"`
}

// Engine description
type EngineDescription struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Language      *string `json:"language"`
	DocumentCount int     `json:"document_count"`
}

type Sorting = stringMap

type SearchGroup struct {
	Field string  `json:"field"`
	Size  int     `json:"size,omitempty"`
	Sort  Sorting `json:"sort,omitempty"`
	// TODO: IDK which type this field must have. Definition is unclear in spec.
	Collapse interface{} `json:"collapse,omitempty"`
}

type FacetType = string

const (
	ValueFacet FacetType = "value"
	RangeFacet FacetType = "range"
)

type Range struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type SearchFilters = m

type Facet struct {
	Type   FacetType `json:"type"`
	Name   string    `json:"name"`
	Sort   Sorting   `json:"sort,omitempty"`
	Size   int       `json:"size,omitempty"`
	Ranges []Range   `json:"ranges,omitempty"`
}

type SearchFacets = map[string][]Facet

type BoostType = string

const (
	ValueBoost      BoostType = "value"
	ProximityBoost  BoostType = "proximity"
	FunctionalBoost BoostType = "functional"
)

type BoostOperation = string

const (
	AddOperation      BoostOperation = "add"
	MultiplyOperation BoostOperation = "multiply"
)

type BoostFunction = string

const (
	LinearFunction      BoostFunction = "linear"
	GaussianFunction    BoostFunction = "gaussian"
	ExponentialFunction BoostFunction = "exponential"
)

type SearchBoost struct {
	Type   BoostType   `json:"type"`
	Value  interface{} `json:"value,omitempty"`
	Factor float32     `json:"factor,omitempty"`
	// Operation for ValueBoost or FunctionalBoost
	Operation BoostOperation `json:"operation,omitempty"`
	// Function for FunctionalBoost or ProximityBoost
	Function BoostFunction `json:"function,omitempty"`
	// Center for ProximityBoost
	Center string `json:"center,omitempty"`
}

type SearchBoosts = map[string]SearchBoost

type SearchAnalytics struct {
	Tags []string `json:"tags"`
}

type FieldWithWeight struct {
	Weight float32 `json:"weight"`
}

type SearchFields = map[string]FieldWithWeight

type RawField struct {
	Size int `json:"size,omitempty"`
}

type SnippetField struct {
	Size     int  `json:"size,omitempty"`
	Fallback bool `json:"fallback,omitempty"`
}

type ResultField struct {
	Raw     *RawField     `json:"raw,omitempty"`
	Snippet *SnippetField `json:"snippet,omitempty"`
}

type ResultFields = map[string]ResultField

// Search query structure
// TODO: query builder
type Query struct {
	// Lucene query
	Query string `json:"query"`
	// Pagination
	Page *Page `json:"page,omitempty"`
	// Sorting
	Sort Sorting `json:"sort,omitempty"`
	// Search grouping
	Group *SearchGroup `json:"group,omitempty"`
	// Search facets
	Facets SearchFacets `json:"facets,omitempty"`
	// Search filters
	Filters SearchFilters `json:"filters,omitempty"`
	// Search boosts
	Boosts SearchBoosts `json:"boosts,omitempty"`
	// Search fields
	SearchFields SearchFields `json:"search_fields,omitempty"`
	// Result fields
	ResultFields ResultFields `json:"result_fields,omitempty"`
	// Analytics
	Analytics *SearchAnalytics `json:"analytics,omitempty"`
}

// Document Search API response
type SearchResponse struct {
	Meta    ResponseMeta `json:"meta"`
	Results []m          `json:"results"`
}

// API Error
type Error struct {
	Message    string   `json:"error"`
	Messages   []string `json:"errors"`
	StatusCode int      `json:"code"`
}

func (e *Error) Error() string {
	if len(e.Messages) > 0 {
		return strings.Join(e.Messages, ", ")
	}

	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("HTTP [%d]", e.StatusCode)
}
