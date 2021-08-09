package tsclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"emperror.dev/errors"
)

// ErrNotSlice is returned by Import if a type other than a slice is given as input.
const ErrNotSlice = errors.Sentinel("slice expected")

// Insert inserts a document into the collection.
// The inserted document is unmarshaled to `out` if it is not nil.
func (c *Client) Insert(collection string, doc interface{}, out interface{}) (err error) {
	resp, err := c.Request("POST", "/collections/"+collection+"/documents", WithJSONBody(doc))
	if err != nil || out == nil {
		return
	}

	return json.Unmarshal(resp, out)
}

// Upsert inserts a document into the collection, updating it if it already exists.
// The inserted document is unmarshaled to `out` if it is not nil.
func (c *Client) Upsert(collection string, doc interface{}, out interface{}) (err error) {
	resp, err := c.Request("POST", "/collections/"+collection+"/documents",
		WithJSONBody(doc),
		WithURLValues(url.Values{"action": {"upsert"}}),
	)
	if err != nil || out == nil {
		return
	}

	return json.Unmarshal(resp, out)
}

type importResponse struct {
	Success bool `json:"success"`
}

// Import imports the documents into the collection.
// It returns an error if s is not a slice.
func (c *Client) Import(collection, action string, s interface{}) (ok []bool, err error) {
	slice := []interface{}{}

	val := reflect.ValueOf(s)

	if val.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	for i := 0; i < val.Len(); i++ {
		slice = append(slice, val.Index(i).Interface())
	}

	return c.ImportSlice(collection, action, slice)
}

// ImportSlice imports a slice of documents into the collection.
func (c *Client) ImportSlice(collection, action string, s []interface{}) (ok []bool, err error) {
	b := new(bytes.Buffer)

	enc := json.NewEncoder(b)

	for _, doc := range s {
		err = enc.Encode(doc)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.Request("POST", "/collections/"+collection+"/documents/import",
		WithBody(b),
		WithHeader(http.Header{
			"Content-Type": {"application/json"},
		}),
	)
	if err != nil {
		return
	}

	dec := json.NewDecoder(bytes.NewReader(resp))

	for dec.More() {
		var r importResponse
		err = dec.Decode(&r)
		if err != nil {
			return
		}

		ok = append(ok, r.Success)
	}
	return
}

// Document retrieves a document from the collection by ID.
// The document is unmarshaled to `out` if it is not nil.
func (c *Client) Document(collection, id string, out interface{}) (string, error) {
	resp, err := c.Request("GET", "/collections/"+collection+"/documents/"+id)
	if err != nil {
		return "", err
	}

	s := struct {
		ID string `json:"id"`
	}{}

	err = json.Unmarshal(resp, &s)
	if err != nil {
		return "", err
	}

	if out == nil {
		return s.ID, nil
	}

	return s.ID, json.Unmarshal(resp, out)
}

// UpdateDocument updates a document in the collection by ID.
// The updated document is unmarshaled to `out` if it is not nil.
func (c *Client) UpdateDocument(collection, id string, doc, out interface{}) error {
	resp, err := c.Request("PATCH", "/collections/"+collection+"/documents/"+id,
		WithJSONBody(doc))
	if err != nil || out == nil {
		return err
	}

	return json.Unmarshal(resp, out)
}

// DeleteDocument deletes a document in the collection.
// The deleted document is unmarshaled to `out` if it is not nil.
func (c *Client) DeleteDocument(collection, id string, out interface{}) error {
	resp, err := c.Request("DELETE", "/collections/"+collection+"/documents/"+id)
	if err != nil || out == nil {
		return err
	}

	return json.Unmarshal(resp, out)
}

// DeleteQuery deletes documents in the collection matching filter.
// Returns the number of deleted documents.
func (c *Client) DeleteQuery(collection, filter string, batchSize int) (deleted int, err error) {
	v := url.Values{"filter_by": {filter}}

	if batchSize != 0 {
		v["batch_size"] = []string{strconv.Itoa(batchSize)}
	}

	resp, err := c.Request("DELETE", "/collections/"+collection+"/documents", WithURLValues(v))
	if err != nil {
		return
	}

	s := struct {
		NumDeleted int `json:"num_deleted"`
	}{}

	err = json.Unmarshal(resp, &s)
	if err != nil {
		return
	}

	return s.NumDeleted, nil
}
