package initialize

import "github.com/pkg/errors"

var (
	ErrFlagsNotSpecified = errors.New("all of available flags has to te specified")
)
