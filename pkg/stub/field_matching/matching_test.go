package field_matching_test

import (
	"encoding/json"
	"testing"

	"github.com/ch629/mockservice/pkg/stub/field_matching"
	"github.com/stretchr/testify/require"
)

func Test_Matcher(t *testing.T) {
	for name, test := range map[string]struct {
		// Input JSON
		json string
		// The type of matcher the JSON should create
		matcherType string
		// Input strings which should be matched
		matchingStrings []string
		// Input strings which should not be matched
		unmatchingStrings []string
	}{
		"EqualTo": {
			json:              `{"equal_to": "foo"}`,
			matcherType:       "equal_to",
			matchingStrings:   []string{"foo"},
			unmatchingStrings: []string{"", "bar", "baz"},
		},
		"StartsWith": {
			json:              `{"starts_with": "foo"}`,
			matcherType:       "starts_with",
			matchingStrings:   []string{"foo", "foobar"},
			unmatchingStrings: []string{"", "bar", "baz", "f", "fo"},
		},
		"EndsWith": {
			json:              `{"ends_with": "foo"}`,
			matcherType:       "ends_with",
			matchingStrings:   []string{"foo", "barfoo"},
			unmatchingStrings: []string{"", "bar", "baz", "f", "fo", "barfo", "bazf"},
		},
		"Contains": {
			json:              `{"contains": "foo"}`,
			matcherType:       "contains",
			matchingStrings:   []string{"foo", "barfoo", "barfoobaz"},
			unmatchingStrings: []string{"", "bar", "baz", "f", "fo", "barfo", "bazf"},
		},
		"Regex": {
			json:              `{"pattern": "(foo){3}"}`,
			matcherType:       "pattern",
			matchingStrings:   []string{"foofoofoo"},
			unmatchingStrings: []string{"", "bar", "baz", "foo", "foofoo", "foofoofo"},
		},
		"And": {
			json:              `{"and": [{"starts_with": "f"}, {"ends_with": "o"}]}`,
			matcherType:       "and",
			matchingStrings:   []string{"foo", "fbarbazo"},
			unmatchingStrings: []string{"", "f", "bar", "baz", "foobar"},
		},
		"Or": {
			json:              `{"or": [{"starts_with": "f"}, {"ends_with": "o"}]}`,
			matcherType:       "or",
			matchingStrings:   []string{"foo", "fbarbazo", "f", "baro", "foobar"},
			unmatchingStrings: []string{"", "bar", "baz"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			matcher, err := field_matching.UnmarshalJSONToFieldMatcher([]byte(test.json))
			require.NoError(t, err)

			require.Equal(t, test.matcherType, matcher.Type())
			for _, matching := range test.matchingStrings {
				require.Truef(t, matcher.Matches(matching), "%s should match", matching)
			}

			for _, unmatching := range test.unmatchingStrings {
				require.Falsef(t, matcher.Matches(unmatching), "%s should not match", unmatching)
			}

			bs, err := json.Marshal(matcher)
			require.NoError(t, err)
			require.JSONEq(t, test.json, string(bs))
		})
	}
}
