package matching

import (
	"encoding/json"
	"fmt"
)

func init() {
	register(func() FieldMatcher { return &orMatcher{} })
}

type orMatcher struct {
	matchers []FieldMatcher
}

func (m orMatcher) Matches(field string) bool {
	for _, m := range m.matchers {
		if m.Matches(field) {
			return true
		}
	}
	return false
}

func OrMatcher(matchers ...FieldMatcher) FieldMatcher {
	return &orMatcher{
		matchers: matchers,
	}
}

func (m orMatcher) MarshalJSON() ([]byte, error) {
	bs, err := json.Marshal(m.matchers)
	if err != nil {
		return nil, fmt.Errorf("marshal to JSON: %w", err)
	}
	return []byte(fmt.Sprintf(`{"or": %s}`, string(bs))), nil
}

func (m *orMatcher) UnmarshalJSON(bs []byte) (err error) {
	var orBlock struct {
		Or []json.RawMessage `json:"or"`
	}
	if err = json.Unmarshal(bs, &orBlock); err != nil {
		return fmt.Errorf("json.Unmarshal to or matcher: %w", err)
	}

	m.matchers = make([]FieldMatcher, len(orBlock.Or))
	// Unmarshal each Matcher inside of the or block
	for idx, orField := range orBlock.Or {
		if m.matchers[idx], err = UnmarshalJSONToFieldMatcher(orField); err != nil {
			return fmt.Errorf("matching.unmarshalJSON: %w", err)
		}
	}
	return
}

func (orMatcher) Type() string {
	return "or"
}
