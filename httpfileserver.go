package httpservefile

import (
	"net/http"
)

// HTTPFileServer define interface for serving file with HTTP handlers.
type HTTPFileServer interface {
	// ServeHTTP respond content to HTTP request `r` with given targetFileName.
	// The defaultFileName will be use instead if targetFileName is empty.
	ServeHTTP(w http.ResponseWriter, r *http.Request, defaultFileName, targetFileName string)

	// Close free used resources.
	Close() (err error)
}
