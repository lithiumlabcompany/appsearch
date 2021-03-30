package appsearch

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-resty/resty/v2"
)

type client struct{ *resty.Client }

func (c *client) Call(ctx context.Context, requestBody, resultPtr interface{}, method, urlFormat string, args ...interface{}) error {
	r, err := c.request(ctx, requestBody, resultPtr).
		Execute(method, fmt.Sprintf(urlFormat, args...))
	if err != nil {
		return err
	}

	if r.IsError() {
		err := r.Error().(*Error)
		err.StatusCode = r.StatusCode()
		// Map error to known api errors for convenience
		if err, ok := apiErrors[err.Error()]; ok {
			return err
		}
		return err
	}

	if resultPtr != nil {
		outElem := reflect.ValueOf(resultPtr).Elem()
		resultElem := reflect.ValueOf(r.Result()).Elem()

		if outElem.Type() != resultElem.Type() {
			return fmt.Errorf("cannot assign result: different types: %s != %s",
				outElem.Type().String(), resultElem.Type().String())
		}

		outElem.Set(resultElem)
	}

	return nil
}

func (c *client) request(ctx context.Context, requestBody interface{}, resultPtr interface{}) *resty.Request {
	req := c.R().
		SetBody(requestBody).
		SetError(&Error{}).
		SetContext(ctx)

	if resultPtr != nil {
		req.SetResult(resultPtr)
	}

	return req
}
