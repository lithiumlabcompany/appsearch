appsearch
=========

[![codecov](https://codecov.io/gh/lithiumlabcompany/appsearch/branch/master/graph/badge.svg?token=qghHHh5Cu6)](https://codecov.io/gh/lithiumlabcompany/appsearch)

Unofficial Experimental AppSearch API client for Go.

[Godoc](https://pkg.go.dev/github.com/lithiumlabcompany/appsearch)
| [Elastic AppSearch API](https://www.elastic.co/guide/en/app-search/current/api-reference.html)

## Features

- Marshal/Unmarshal structure to schema
- Engine API [Godoc](https://pkg.go.dev/github.com/lithiumlabcompany/appsearch#EngineAPI)
  | [ElasticSearch Reference](https://www.elastic.co/guide/en/app-search/current/engines.html)
- Schema API [Godoc](https://pkg.go.dev/github.com/lithiumlabcompany/appsearch#SchemaAPI)
  | [ElasticSearch Reference](https://www.elastic.co/guide/en/app-search/current/schema.html)
- Document
  API [Godoc](https://pkg.go.dev/github.com/lithiumlabcompany/appsearch#DocumentAPI)
  | [ElasticSearch Reference](https://www.elastic.co/guide/en/app-search/current/documents.html)

## TODO

- Derive schema from struct with tags
- Implement complete set of Elastic App Search API's

## Quickstart

```go
package main

import (
	"context"
	"github.com/lithiumlabcompany/appsearch"
	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

type Civilization struct {
	Name        string
	Rating      float32
	Description string
}

func main() {
	client, _ := appsearch.Open("https://private-key@endpoint.ent-search.cloud.es.io")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engineName := "civilizations"
	schemaDefinition := appsearch.SchemaDefinition{
		"name":        "text",
		"rating":      "number",
		"description": "text",
	}
	// Engine will be created if it doesn't exist and schema will be updated
	client.EnsureEngine(ctx, appsearch.CreateEngineRequest{
		Name:     "civilizations",
		Language: "en",
	}, schemaDefinition)

	document, _ := schema.ToMap(Civilization{
		Name:        "Babylonian",
		Rating:      5212.2,
		Description: "Technological and scientific",
	}, schemaDefinition)

	client.UpdateDocuments(ctx, "civilizations", []schema.Map{document})

	search, _ := client.SearchDocuments(ctx, engineName, appsearch.Query{
		Query: "scientific",
	})

	var results []Civilization
	_ = schema.UnmarshalResults(search.Results, &results)

	println(results[0])
}
```
