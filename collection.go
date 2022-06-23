package tsclient

import "encoding/json"

// Collection ...
type Collection struct {
	Name                string  `json:"name"`
	NumDocuments        int     `json:"num_documents,omitempty"`
	Fields              []Field `json:"fields"`
	DefaultSortingField string  `json:"default_sorting_field,omitempty"`
}

// Field is a field of a collection.
type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Facet bool   `json:"facet"`
	Index bool   `json:"index"`
	Infix bool   `json:"infix"`
}

// Collection gets a collection by name.
func (c *Client) Collection(name string) (col Collection, err error) {
	resp, err := c.Request("GET", "/collections/"+name)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &col)
	return
}

// Collections returns a summary of all your collections.
// The collections are returned sorted by creation date,
// with the most recent collections appearing first.
func (c *Client) Collections() (cols []Collection, err error) {
	resp, err := c.Request("GET", "/collections")
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &cols)
	return
}

// DeleteCollection permanently drops a collection. This action cannot be undone.
// For large collections, this might have an impact on read latencies.
func (c *Client) DeleteCollection(name string) (col Collection, err error) {
	resp, err := c.Request("DELETE", "/collections/"+name)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &col)
	return
}

// CreateFieldData is the field data passed to CreateCollection.
// The Index field is inverted here to avoid needing a bool pointer.
type CreateFieldData struct {
	Name  string
	Type  string
	Facet bool
	// false = index the field, true = don't index the field
	NoIndex bool
	Infix   bool
}

// CreateCollection creates a collection. defaultSortingField is optional and may be left empty.
func (c *Client) CreateCollection(name, defaultSortingField string, fields []CreateFieldData) (col Collection, err error) {
	fs := []Field{}
	for _, f := range fields {
		fs = append(fs, Field{
			Name:  f.Name,
			Type:  f.Type,
			Facet: f.Facet,
			Index: !f.NoIndex,
			Infix: f.Infix,
		})
	}

	resp, err := c.Request("POST", "/collections", WithJSONBody(&Collection{
		Name:                name,
		DefaultSortingField: defaultSortingField,
		Fields:              fs,
	}))
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &col)
	return
}
