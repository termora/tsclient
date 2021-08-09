package tsclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// RequestOption is an optional request option.
type RequestOption func(*http.Request) error

// WithHeader adds headers to the request.
func WithHeader(header http.Header) func(*http.Request) error {
	return func(req *http.Request) error {
		for k, v := range header {
			req.Header[k] = append(req.Header[k], v...)
		}

		return nil
	}
}

// WithURLValues adds query parameters to the request.
func WithURLValues(values url.Values) func(*http.Request) error {
	return func(req *http.Request) error {
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// WithBody adds a body to the request.
func WithBody(r io.Reader) func(*http.Request) error {
	return func(req *http.Request) error {
		rc, ok := r.(io.ReadCloser)
		if !ok && r != nil {
			rc = io.NopCloser(r)
		}

		req.Body = rc

		return nil
	}
}

// WithJSONBody adds a JSON body to the request. Marshaling errors are ignored.
func WithJSONBody(v interface{}) func(*http.Request) error {
	if v == nil {
		return func(*http.Request) error { return nil }
	}

	return func(req *http.Request) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}
}
