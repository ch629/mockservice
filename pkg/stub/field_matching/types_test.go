package field_matching

func init() {
	register(func() FieldMatcher { return &TrueMatcher{} })
	register(func() FieldMatcher { return &FalseMatcher{} })
}

type TrueMatcher struct{}

func (m *TrueMatcher) Matches(field string) bool {
	return true
}

func (m TrueMatcher) MarshalJSON() ([]byte, error) {
	return []byte(`{"true": true}`), nil
}

func (m *TrueMatcher) UnmarshalJSON(bs []byte) error {
	return nil
}

func (TrueMatcher) Type() string {
	return "true"
}

type FalseMatcher struct{}

func (m *FalseMatcher) Matches(field string) bool {
	return false
}

func (m FalseMatcher) MarshalJSON() ([]byte, error) {
	return []byte(`{"false": false}`), nil
}

func (m *FalseMatcher) UnmarshalJSON(bs []byte) error {
	return nil
}

func (FalseMatcher) Type() string {
	return "false"
}
