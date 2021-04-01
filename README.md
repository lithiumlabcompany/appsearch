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
)

func main() {
	client, _ := appsearch.Open("https://private-key@endpoint.ent-search.cloud.es.io")

	// Engine will be created if it doesn't exist and schema will be updated
	client.EnsureEngine(ctx, appsearch.CreateEngineRequest{
		Name:     "civilizations",
		Language: "en",
	}, appsearch.SchemaDefinition{
		"name":        "text",
		"rating":      "number",
		"description": "text",
	})

	client.UpdateDocuments(ctx, "civilizations", []map[string]interface{}{
		{"name": "Babylonian", "rating": 5212.2, "description": "Technological and scientific"},
	})

	search, _ := client.SearchDocuments(ctx, "civilizations", appsearch.Query{
		Query: "science",
	})

	println(search.Results[0])

	/*
		{
		  "_meta": {
		    "score": 1396363.1
		  },
		  "name": {
		    "raw": "Babylonian"
		  },
		  "description": {
		    "raw": "Technological and scientific"
		  },
		  "rating": {
		    "raw": 5212.2
		  },
		  "id": {
		    "raw": "park_everglades"
		  }
		}
	*/
}
```
