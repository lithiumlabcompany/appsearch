package mock

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/lithiumlabcompany/appsearch"
	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

type mock struct {
	Engines map[string]appsearch.EngineDescription
	Schemas map[string]schema.Definition

	Implementation map[string]interface{}
}

// Create mock APIClient
func Mock(args ...interface{}) *mock {
	m := &mock{
		Engines:        map[string]appsearch.EngineDescription{},
		Schemas:        map[string]schema.Definition{},
		Implementation: map[string]interface{}{},
	}
	for _, v := range args {
		switch v := v.(type) {
		case map[string]appsearch.EngineDescription:
			m.Engines = v
		case map[string]schema.Definition:
			m.Schemas = v
		case map[string]interface{}:
			m.Implementation = v
		default:
			panic(fmt.Errorf("accepted params for Mock() are only updates on Engine, Schemas or Implementations fields"))
		}
	}
	return m
}

func (m *mock) impl(args []interface{}, resultPointers []interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	addr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	method := addr[len(addr)-1]

	implFunc, ok := m.Implementation[method]
	if !ok {
		panic(fmt.Errorf("implement me (%s)", method))
	}

	methodType := reflect.ValueOf(implFunc).Type()
	resultValues := reflect.ValueOf(implFunc).Call(callValues(methodType, args))
	assignPointerValues(resultValues, resultPointers)
}

func assignPointerValues(values []reflect.Value, pointers []interface{}) {
	for i, value := range values {
		pointer := pointers[i]

		reflect.ValueOf(pointer).Elem().Set(value)
	}
}

func callValues(method reflect.Type, interfaces []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(interfaces))
	for i := range interfaces {
		valueType := method.In(i)
		if interfaces[i] == nil {
			values[i] = reflect.New(valueType).Elem()
		} else {
			values[i] = reflect.ValueOf(interfaces[i])
		}
	}
	return values
}

func sliceAsInterfaces(slice interface{}) []interface{} {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		interfaces := make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			interfaces[i] = s.Index(i).Interface()
		}

		return interfaces
	default:
		panic(fmt.Errorf("value of %v is not reflect.Slice", slice))
	}
}

func interfacesOf(arg ...interface{}) []interface{} {
	return arg
}
