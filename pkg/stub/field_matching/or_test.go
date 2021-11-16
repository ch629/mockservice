package field_matching_test

import (
	"encoding/json"
	"testing"

	"github.com/ch629/mockservice/pkg/stub/field_matching"
	"github.com/stretchr/testify/require"
)

func Test_or_Matches(t *testing.T) {
	for name, test := range map[string]struct {
		matchers      []field_matching.FieldMatcher
		expectedValue bool
	}{
		"No matchers": {
			matchers:      []field_matching.FieldMatcher{},
			expectedValue: false,
		},
		"True matcher": {
			matchers:      []field_matching.FieldMatcher{&field_matching.TrueMatcher{}},
			expectedValue: true,
		},
		"False matcher": {
			matchers:      []field_matching.FieldMatcher{&field_matching.FalseMatcher{}},
			expectedValue: false,
		},
		"Many matchers": {
			matchers: []field_matching.FieldMatcher{
				&field_matching.FalseMatcher{},
				&field_matching.FalseMatcher{},
				&field_matching.FalseMatcher{},
				&field_matching.FalseMatcher{},
				&field_matching.TrueMatcher{},
			},
			expectedValue: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := field_matching.OrMatcher(test.matchers...)
			require.Equal(t, test.expectedValue, matcher.Matches(""))
		})
	}
}

func Test_or_MarshalUnmarshal(t *testing.T) {
	for name, test := range map[string]struct {
		expectedJSON string
		matchers     []field_matching.FieldMatcher
	}{
		"Multiple matchers": {
			expectedJSON: `
			{
				"or": [
					{"true": true},
					{"false": false},
					{"false": false}
				]
			}`,
			matchers: []field_matching.FieldMatcher{&field_matching.TrueMatcher{}, &field_matching.FalseMatcher{}, &field_matching.FalseMatcher{}},
		},
		"Empty matchers": {
			expectedJSON: `
			{ 
				"or": []
			}`,
			matchers: []field_matching.FieldMatcher{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := field_matching.OrMatcher(test.matchers...)
			bs, err := json.Marshal(matcher)
			require.NoError(t, err)
			require.JSONEq(t, test.expectedJSON, string(bs))

			newMatcher := field_matching.OrMatcher()
			err = newMatcher.UnmarshalJSON(bs)
			require.NoError(t, err)
			require.Equal(t, matcher, newMatcher)
		})
	}
}
