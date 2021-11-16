package field_matching

import "encoding/json"

type FieldMatcher interface {
	json.Unmarshaler
	json.Marshaler
	Matches(field string) bool
	Type() string
}
