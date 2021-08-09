// Package jsonutil implements types and functions to make it easier to handle the JSON returned by Typesense.
package jsonutil

import "encoding/json"

// Raw stores raw JSON data, to be unmarshaled at a later time.
type Raw []byte

// UnmarshalJSON sets *m to a copy of data.
func (m *Raw) UnmarshalJSON(data []byte) error {
	*m = append((*m)[0:0], data...)
	return nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m Raw) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalTo unmarshals m into v.
func (m Raw) UnmarshalTo(v interface{}) error {
	if len(m) == 0 {
		return nil
	}

	return json.Unmarshal(m, v)
}

// StringPointer returns a pointer to s.
func StringPointer(s string) *string { return &s }

// IntPointer returns a pointer to i.
func IntPointer(i int) *int { return &i }

// BoolPointer returns a pointer to b.
func BoolPointer(b bool) *bool { return &b }
