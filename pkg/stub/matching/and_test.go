package matching_test

import (
	"encoding/json"
	"testing"

	"github.com/ch629/mockservice/pkg/stub/matching"
	"github.com/stretchr/testify/require"
)

func Test_and_Matches(t *testing.T) {
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
		"All true": {
			matchers: []matching.FieldMatcher{
				&matching.TrueMatcher{},
				&matching.TrueMatcher{},
				&matching.TrueMatcher{},
				&matching.TrueMatcher{},
				&matching.TrueMatcher{},
			},
			expectedValue: true,
		},
		"Many matchers": {
			matchers: []matching.FieldMatcher{
				&matching.TrueMatcher{},
				&matching.TrueMatcher{},
				&matching.FalseMatcher{},
				&matching.FalseMatcher{},
				&matching.FalseMatcher{},
			},
			expectedValue: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := matching.AndMatcher(test.matchers...)
			require.Equal(t, test.expectedValue, matcher.Matches(""))
		})
	}
}

func Test_and_MarshalUnmarshal(t *testing.T) {
	for name, test := range map[string]struct {
		expectedJSON string
		matchers     []matching.FieldMatcher
	}{
		"Multiple matchers": {
			expectedJSON: `
			{
				"and": [
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
				"and": []
			}`,
			matchers: []matching.FieldMatcher{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher := matching.AndMatcher(test.matchers...)
			bs, err := json.Marshal(matcher)
			require.NoError(t, err)
			require.JSONEq(t, test.expectedJSON, string(bs))

			newMatcher := matching.AndMatcher()
			err = newMatcher.UnmarshalJSON(bs)
			require.NoError(t, err)
			require.Equal(t, matcher, newMatcher)
		})
	}
}
