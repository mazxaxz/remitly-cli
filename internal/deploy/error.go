package deploy

import "github.com/pkg/errors"

var (
	ErrPathVariableNotSet        = errors.New("REMITLY_PATH environment variable not set")
	ErrProfileVariableNotSet     = errors.New("REMITLY_PROFILE environment variable not set")
	ErrInvalidContextsFileSyntax = errors.New("invalid $REMITLY_PATH/*.yml file syntax")
	ErrProfileNotFound           = errors.New("profile $REMITLY_PROFILE was not found inside $REMITLY_PATH/*.yml file")
)
