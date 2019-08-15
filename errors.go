package httpservefile

import (
	"errors"
)

// ErrEmptyContentStoragePath indicate path to content storage is empty.
var ErrEmptyContentStoragePath = errors.New("content storage path is not given")
