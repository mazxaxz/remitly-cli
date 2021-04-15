package deploy

import "github.com/pkg/errors"

var (
	ErrPathVariableNotFound = errors.New("REMITLY_PATH environment variable not found")
)
