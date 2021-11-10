package matching

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func init() {
	for _, name := range []string{
		"pattern",
		"equal_to",
		"contains",
		"starts_with",
		"ends_with",
	} {
		register(name, func() FieldMatcher { return &simpleMatcher{} })
	}
}

type simpleMatcher struct {
	matchPredicate func(str string) bool
	name           string
	field          string
}

func (m simpleMatcher) Matches(field string) bool {
	return m.matchPredicate(field)
}

func (m simpleMatcher) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s": "%s"}`, m.name, m.field)), nil
}

func (m simpleMatcher) Type() string {
	return m.name
}

func (m *simpleMatcher) UnmarshalJSON(bs []byte) error {
	var catchAll struct {
		Pattern    string `json:"pattern"`
		EqualTo    string `json:"equal_to"`
		Contains   string `json:"contains"`
		StartsWith string `json:"starts_with"`
		EndsWith   string `json:"ends_with"`
	}
	if err := json.Unmarshal(bs, &catchAll); err != nil {
		return fmt.Errorf("json.Unmarshal to catchAll: %w", err)
	}
	switch {
	case catchAll.Pattern != "":
		reg, err := RegexMatcher(catchAll.Pattern)
		if err != nil {
			return fmt.Errorf("creating RegexMatcher: %w", err)
		}
		*m = *reg.(*simpleMatcher)
	case catchAll.EqualTo != "":
		*m = *EqualToMatcher(catchAll.EqualTo).(*simpleMatcher)
	case catchAll.Contains != "":
		*m = *ContainsWithMatcher(catchAll.Contains).(*simpleMatcher)
	case catchAll.StartsWith != "":
		*m = *StartsWithMatcher(catchAll.StartsWith).(*simpleMatcher)
	case catchAll.EndsWith != "":
		*m = *EndsWithMatcher(catchAll.EndsWith).(*simpleMatcher)
	default:
		return errors.New("unknown key")
	}
	return nil
}

// TODO: Pass regexp.Regexp instead?
func RegexMatcher(pattern string) (FieldMatcher, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}

	return &simpleMatcher{
		matchPredicate: func(str string) bool {
			return regex.MatchString(str)
		},
		name:  "pattern",
		field: pattern,
	}, nil
}

func EqualToMatcher(match string) FieldMatcher {
	return &simpleMatcher{
		matchPredicate: func(str string) bool {
			return str == match
		},
		name:  "equal_to",
		field: match,
	}
}

func ContainsWithMatcher(contains string) FieldMatcher {
	return &simpleMatcher{
		matchPredicate: func(str string) bool {
			return strings.Contains(str, contains)
		},
		name:  "contains",
		field: contains,
	}
}

func StartsWithMatcher(startsWith string) FieldMatcher {
	return &simpleMatcher{
		matchPredicate: func(str string) bool {
			return strings.HasPrefix(str, startsWith)
		},
		name:  "starts_with",
		field: startsWith,
	}
}

func EndsWithMatcher(endsWith string) FieldMatcher {
	return &simpleMatcher{
		matchPredicate: func(str string) bool {
			return strings.HasSuffix(str, endsWith)
		},
		name:  "ends_with",
		field: endsWith,
	}
}
