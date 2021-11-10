package stub

import (
	"encoding/json"
	"fmt"

	"github.com/ch629/mockservice/pkg/domain"
	"github.com/ch629/mockservice/pkg/stub/matching"
	"go.uber.org/zap"
)

type (
	RequestMatcher interface {
		Matches(req domain.Request) bool
	}

	loggedMatcher struct {
		log     *zap.Logger
		matcher RequestMatcher
	}

	pathMatcher struct {
		fieldMatcher matching.FieldMatcher
	}

	queryParameterMatcher struct {
		fieldMatcher matching.FieldMatcher
		key          string
	}
)

func NewPathMatcher(fieldMatcher matching.FieldMatcher) RequestMatcher {
	return &pathMatcher{
		fieldMatcher: fieldMatcher,
	}
}

func NewQueryParamterMatcher(key string, matcher matching.FieldMatcher) RequestMatcher {
	return &queryParameterMatcher{
		key:          key,
		fieldMatcher: matcher,
	}
}

func NewLoggedMatcher(log *zap.Logger, matcher RequestMatcher) RequestMatcher {
	return &loggedMatcher{
		log:     log,
		matcher: matcher,
	}
}

func (q queryParameterMatcher) MarshalJSON() ([]byte, error) {
	bs, err := json.Marshal(q.fieldMatcher)
	if err != nil {
		return nil, fmt.Errorf("marshalling field matcher: %w", err)
	}
	return []byte(fmt.Sprintf(`{"query_paramter": %s}`, string(bs))), nil
}

func (l loggedMatcher) Matches(req domain.Request) bool {
	match := l.matcher.Matches(req)
	l.log.Debug("checking match",
		zap.String("path", req.Path),
		zap.Bool("matched", match),
		zap.Any("matcher", l.matcher),
	)
	return match
}

func (l loggedMatcher) MarshalJSON() ([]byte, error) {
	//nolint
	return json.Marshal(l.matcher)
}

func (p pathMatcher) Matches(req domain.Request) bool {
	return p.fieldMatcher.Matches(req.Path)
}

func (p pathMatcher) MarshalJSON() ([]byte, error) {
	bs, err := json.Marshal(p.fieldMatcher)
	if err != nil {
		return nil, fmt.Errorf("marshalling field matcher: %w", err)
	}
	return []byte(fmt.Sprintf(`{"path": %s}`, string(bs))), nil
}

func (q queryParameterMatcher) Matches(req domain.Request) bool {
	for _, param := range req.QueryParameters[q.key] {
		if q.fieldMatcher.Matches(param) {
			return true
		}
	}
	return false
}
