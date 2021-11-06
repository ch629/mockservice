package stub

import (
	"fmt"

	"github.com/ch629/mockservice/pkg/domain"
	"go.uber.org/zap"
)

type (
	RequestMatcher interface {
		// TODO: Stringer or json marshaller?
		fmt.Stringer
		Matches(req domain.Request) bool
	}

	loggedMatcher struct {
		log     *zap.Logger
		matcher RequestMatcher
	}

	pathMatcher struct {
		path string
	}

	// TODO: valueEquals needs to be some internal struct to match on query param (equals, startsWith, etc)
	queryParameterMatcher struct {
		key         string
		valueEquals string
	}
)

func NewPathMatcher(path string) RequestMatcher {
	return &pathMatcher{
		path: path,
	}
}

func NewLoggedMatcher(log *zap.Logger, matcher RequestMatcher) RequestMatcher {
	return &loggedMatcher{
		log:     log,
		matcher: matcher,
	}
}

func NewQueryParamterMatcher(key, valueEquals string) RequestMatcher {
	return &queryParameterMatcher{
		key:         key,
		valueEquals: valueEquals,
	}
}

func (l loggedMatcher) Matches(req domain.Request) bool {
	match := l.matcher.Matches(req)
	l.log.Debug("checking match",
		zap.String("path", req.Path),
		zap.Bool("matched", match),
		zap.Stringer("matcher", l.matcher),
	)
	return match
}

func (l loggedMatcher) String() string {
	return l.matcher.String()
}

func (p pathMatcher) Matches(req domain.Request) bool {
	return req.Path == p.path
}

func (p pathMatcher) String() string {
	return fmt.Sprintf(`{"path": "%s"}`, p.path)
}

func (q queryParameterMatcher) Matches(req domain.Request) bool {
	for _, param := range req.QueryParameters[q.key] {
		if param == q.valueEquals {
			return true
		}
	}
	return false
}

func (q queryParameterMatcher) String() string {
	return fmt.Sprintf(`{"queryParameter": {"key": "%s", "equals": "%s"}}`, q.key, q.valueEquals)
}
