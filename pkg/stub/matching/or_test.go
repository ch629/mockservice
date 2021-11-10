package matching_test

import (
	"encoding/json"
	"testing"

	"github.com/ch629/mockservice/pkg/stub/matching"
	"github.com/stretchr/testify/require"
)

func Test_or_Matches(t *testing.T) {
	for name, test := range map[string]struct {
		matchers      []matching.FieldMatcher
		expectedValue bool
	}{
		"No matchers": {
			matchers:      []matching.FieldMatcher{},
			expectedValue: false,
		},
		"True matcher": {
			matchers:      []matching.FieldMatcher{&matching.TrueMatcher{}},
			expectedValue: true,
		},
		"False matcher": {
			matchers:      []matching.FieldMatcher{&matching.FalseMatcher{}},
			expectedValue: false,
		},
		"Many matchers": {
			matchers: []matching.FieldMatcher{
				&matching.FalseMatcher{},
				&matching.FalseMatcher{},
				&matching.FalseMatcher{},
				&matching.FalseMatcher{},
				&matching.TrueMatcher{},
			},
			expectedValue: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := matching.OrMatcher(test.matchers...)
			require.Equal(t, test.expectedValue, matcher.Matches(""))
		})
	}
}

func Test_or_MarshalUnmarshal(t *testing.T) {
	for name, test := range map[string]struct {
		expectedJSON string
		matchers     []matching.FieldMatcher
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
			matchers: []matching.FieldMatcher{&matching.TrueMatcher{}, &matching.FalseMatcher{}, &matching.FalseMatcher{}},
		},
		"Empty matchers": {
			expectedJSON: `
			{ 
				"or": []
			}`,
			matchers: []matching.FieldMatcher{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := matching.OrMatcher(test.matchers...)
			bs, err := json.Marshal(matcher)
			require.NoError(t, err)
			require.JSONEq(t, test.expectedJSON, string(bs))

			newMatcher := matching.OrMatcher()
			err = newMatcher.UnmarshalJSON(bs)
			require.NoError(t, err)
			require.Equal(t, matcher, newMatcher)
		})
	}
}
