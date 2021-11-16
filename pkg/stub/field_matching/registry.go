package field_matching

import (
	"encoding/json"
	"errors"
	"fmt"
)

var typeRegistry = make(map[string]func() FieldMatcher)

func register(constructor func() FieldMatcher) {
	registerNamed(constructor().Type(), constructor)
}

func registerNamed(name string, constructor func() FieldMatcher) {
	if _, ok := typeRegistry[name]; ok {
		panic(fmt.Errorf("type '%s' already registered", name))
	}
	typeRegistry[name] = constructor
}

func UnmarshalJSONToFieldMatcher(bs []byte) (FieldMatcher, error) {
	// There should only be one value in this JSON, grab the name
	name, err := getFirstKey(bs)
	if err != nil {
		return nil, fmt.Errorf("getFirstKey: %w", err)
	}
	matcherConstructor, ok := typeRegistry[name]
	if !ok {
		return nil, fmt.Errorf("no matcher registered under '%s'", name)
	}

	matcher := matcherConstructor()
	if err := json.Unmarshal(bs, &matcher); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return matcher, nil
}

// getFirstKey returns the value of the first key in JSON, or an error if there is none
func getFirstKey(msg json.RawMessage) (string, error) {
	var mp map[string]json.RawMessage
	if err := json.Unmarshal(msg, &mp); err != nil {
		return "", fmt.Errorf("json.Unmarshal to map: %w", err)
	}
	for key := range mp {
		return key, nil
	}
	return "", errors.New("no key found in JSON")
}
