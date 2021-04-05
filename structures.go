package appsearch

import (
	"fmt"
	"strings"

	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

// Page options
type Page struct {
	Page int `json:"current,omitempty"`
	Size int `json:"size,omitempty"`
}

// UpdateResponse Response for Patch or Update operations
type UpdateResponse struct {
	// Updated document ID
	ID string `json:"id"`
	// List of errors
	Errors []string `json:"errors"`
}

// DeleteResponse Response for Patch or Update operations
type DeleteResponse struct {
	// Deleted document ID
	ID string `json:"id"`
	// Was document deleted successfully
	Deleted bool `json:"deleted"`
	// List of errors
	Errors []string `json:"errors"`
}

// PaginationMeta Pagination metadata included in paged responses
type PaginationMeta struct {
	PageSize     int `json:"size"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current"`
	TotalResults int `json:"total_results"`
}

// ResponseMeta Response metadata included in some responses
type ResponseMeta struct {
	Page      PaginationMeta `json:"page"`
	Alerts    []string       `json:"alerts"`
	Warnings  []string       `json:"warnings"`
	RequestID string         `json:"request_id"`
}

// EngineResponse ListEngines response
type EngineResponse struct {
	Meta    ResponseMeta        `json:"meta"`
	Results []EngineDescription `json:"results"`
}

// CreateEngineRequest Request for CreateEngine
type CreateEngineRequest struct {
	Name     string `json:"name"`
	Language string `json:"language,omitempty"`
}

// EngineDescription Engine description
type EngineDescription struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Language      *string `json:"language"`
	DocumentCount int     `json:"document_count"`
}

// Sorting options
type Sorting = map[string]string

// SearchGroup Search group
type SearchGroup struct {
	Field    string  `json:"field"`
	Size     int     `json:"size,omitempty"`
	Sort     Sorting `json:"sort,omitempty"`
	Collapse bool    `json:"collapse,omitempty"`
}

// FacetType sFacet type
type FacetType = string

const (
	// ValueFacet Facet type
	ValueFacet FacetType = "value"
	// RangeFacet Facet type
	RangeFacet FacetType = "range"
)

// Range Search Range
type Range struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// SearchFilters
type SearchFilters = schema.Map

// Facet
type Facet struct {
	Type   FacetType `json:"type"`
	Name   string    `json:"name"`
	Sort   Sorting   `json:"sort,omitempty"`
	Size   int       `json:"size,omitempty"`
	Ranges []Range   `json:"ranges,omitempty"`
}

// SearchFacets Search facets
type SearchFacets = map[string][]Facet

// Describe single facet data in FacetResult
type FacetData struct {
	Value string `json:"value"`
	Count int    `json:"count"`
	// Can be date string, or number
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// Describe single facet result in FacetResultMap
type FacetResult struct {
	Type FacetType   `json:"type"`
	Name string      `json:"name"`
	Data []FacetData `json:"data"`
}

// Structure to describe facets in response
type FacetResultMap = map[string][]FacetResult

// BoostType Boost type
type BoostType = string

const (
	// ValueBoost Boost type
	ValueBoost BoostType = "value"
	// ProximityBoost Boost type
	ProximityBoost BoostType = "proximity"
	// FunctionalBoost Boost type
	FunctionalBoost BoostType = "functional"
)

// BoostOperation Boost operation
type BoostOperation = string

const (
	// AddOperation Boost operation
	AddOperation BoostOperation = "add"
	// MultiplyOperation Boost operation
	MultiplyOperation BoostOperation = "multiply"
)

// BoostFunction Boost function
type BoostFunction = string

const (
	// Boost function
	LinearFunction BoostFunction = "linear"
	// Boost function
	GaussianFunction BoostFunction = "gaussian"
	// Boost function
	ExponentialFunction BoostFunction = "exponential"
)

// Search boosts
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

// Search boosts
type SearchBoosts = map[string]SearchBoost

// Search analytics
type SearchAnalytics struct {
	Tags []string `json:"tags"`
}

// Field with weight
type FieldWithWeight struct {
	Weight float32 `json:"weight"`
}

// Specify search felds
type SearchFields = map[string]FieldWithWeight

// Specify field as raw
type RawField struct {
	Size int `json:"size,omitempty"`
}

// Specify field as snippet
type SnippetField struct {
	Size     int  `json:"size,omitempty"`
	Fallback bool `json:"fallback,omitempty"`
}

// Specify what field must look like in output
type ResultField struct {
	Raw     *RawField     `json:"raw,omitempty"`
	Snippet *SnippetField `json:"snippet,omitempty"`
}

// Map of result field specifications
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
	Meta    ResponseMeta   `json:"meta"`
	Facets  FacetResultMap `json:"facets"`
	Results []schema.Map   `json:"results"`
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
