package matching

import (
	"encoding/json"
	"fmt"
)

func init() {
	register(func() FieldMatcher { return &andMatcher{} })
}

type andMatcher struct {
	matchers []FieldMatcher
}

func (m andMatcher) Matches(field string) bool {
	if len(m.matchers) == 0 {
		return false
	}

	for _, m := range m.matchers {
		if !m.Matches(field) {
			return false
		}
	}
	return true
}

func AndMatcher(matchers ...FieldMatcher) FieldMatcher {
	return &andMatcher{
		matchers: matchers,
	}
}

func (m andMatcher) MarshalJSON() ([]byte, error) {
	bs, err := json.Marshal(m.matchers)
	if err != nil {
		return nil, fmt.Errorf("marshal to JSON: %w", err)
	}
	return []byte(fmt.Sprintf(`{"and": %s}`, string(bs))), nil
}

func (m *andMatcher) UnmarshalJSON(bs []byte) (err error) {
	var andBlock struct {
		And []json.RawMessage `json:"and"`
	}
	if err = json.Unmarshal(bs, &andBlock); err != nil {
		return fmt.Errorf("json.Unmarshal to or matcher: %w", err)
	}

	m.matchers = make([]FieldMatcher, len(andBlock.And))
	// Unmarshal each Matcher inside of the and block
	for idx, andField := range andBlock.And {
		if m.matchers[idx], err = UnmarshalJSONToFieldMatcher(andField); err != nil {
			return fmt.Errorf("matching.unmarshalJSON: %w", err)
		}
	}
	return
}

func (andMatcher) Type() string {
	return "and"
}
