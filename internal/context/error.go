package context

import "github.com/pkg/errors"

var (
	ErrEmptyUseFlag = errors.New("flag --use must contain any value")
)
