package tsclient

import (
	"fmt"
	"io"
	"net/http"

	"emperror.dev/errors"
)

// Errors returned by Request
const (
	ErrBadRequest    = errors.Sentinel("400 bad request")
	ErrUnauthorized  = errors.Sentinel("401 unauthorized")
	ErrNotFound      = errors.Sentinel("404 not found")
	ErrAlreadyExists = errors.Sentinel("409 resource already exists")
	ErrUnprocessable = errors.Sentinel("422 unprocessable entity")
	ErrUnavailable   = errors.Sentinel("503 service unavailable")
)

type apiError int

func (e apiError) Error() string {
	return fmt.Sprintf("%v %v", int(e), http.StatusText(int(e)))
}

// Request makes a request returning a JSON body.
func (c *Client) Request(method, endpoint string, opts ...RequestOption) (response []byte, err error) {
	c.Debug("Request to %v (%v)", endpoint, method)

	req, err := http.NewRequest(method, c.baseURL+endpoint, nil)
	if err != nil {
		return
	}

	for _, opt := range opts {
		err = opt(req)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header["X-TYPESENSE-API-KEY"] = []string{c.apiKey}

	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			c.Debug("error closing response body: %v", err2)
		}
	}()

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent, http.StatusCreated:
	case http.StatusBadRequest:
		return
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusConflict:
		return nil, ErrAlreadyExists
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessable
	case http.StatusServiceUnavailable:
		return nil, ErrUnavailable
	default:
		return nil, apiError(resp.StatusCode)
	}

	return response, err
}
