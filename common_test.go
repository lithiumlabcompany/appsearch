package appsearch

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func createRandomEngine(c APIClient) EngineDescription {
	engine, err := c.CreateEngine(context.TODO(), CreateEngineRequest{
		Name: strings.ToLower(uuid.New().String()),
	})
	exit(err)
	return engine
}

func deleteEngine(c APIClient, engine interface{}) {
	switch e := engine.(type) {
	case EngineDescription:
		exit(c.DeleteEngine(context.TODO(), e.Name))
	case string:
		exit(c.DeleteEngine(context.TODO(), e))
	default:
		panic(e)
	}
}

func exit(err error) {
	if err != nil {
		panic(err)
	}
}
