package stub

import (
	"fmt"

	"go.uber.org/zap"
)

type (
	RequestMatcher interface {
		// TODO: Stringer or json marshaller?
		fmt.Stringer
		Matches(req Request) bool
	}

	loggedMatcher struct {
		logger  *zap.Logger
		matcher RequestMatcher
	}

	pathMatcher struct {
		path string
	}
)

func NewPathMatcher(path string) RequestMatcher {
	return &pathMatcher{
		path: path,
	}
}

func NewLoggedMatcher(logger *zap.Logger, matcher RequestMatcher) RequestMatcher {
	return &loggedMatcher{
		logger:  logger,
		matcher: matcher,
	}
}

func (l loggedMatcher) Matches(req Request) bool {
	match := l.matcher.Matches(req)
	l.logger.Debug("checking match", zap.String("path", req.Path), zap.Bool("matched", match), zap.Stringer("matcher", l.matcher))
	return match
}

func (l loggedMatcher) String() string {
	return fmt.Sprintf(`{"wrapped": %s}`, l.matcher.String())
}

func (p pathMatcher) Matches(req Request) bool {
	return req.Path == p.path
}

func (p pathMatcher) String() string {
	return fmt.Sprintf(`{"path": "%s"}`, p.path)
}
