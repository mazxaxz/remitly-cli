package deploy

import "github.com/pkg/errors"

var (
	ErrPathVariableNotSet          = errors.New("REMITLY_PATH environment variable not set")
	ErrProfileVariableNotSet       = errors.New("REMITLY_PROFILE environment variable not set")
	ErrInvalidContextsFileSyntax   = errors.New("invalid $REMITLY_PATH/*.yml file syntax")
	ErrProfileNotFound             = errors.New("profile $REMITLY_PROFILE was not found inside $REMITLY_PATH/*.yml file")
	ErrReplicaCountMustBeAboveZero = errors.New("value of --replica-count flag must be above zero")
	ErrFailedDeployment            = errors.New("deployment has failed")
	ErrVersionAlreadyDeployed      = errors.New("given app version has been already deployed before")
)
