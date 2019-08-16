package httpservefile

import (
	"net/http"
)

// sanitizeURLPathPrefix make sure result will start and end with slash (`/`) character.
func sanitizeURLPathPrefix(urlPathPrefix string) string {
	if len(urlPathPrefix) == 0 {
		urlPathPrefix = "/"
	}
	if urlPathPrefix[0] != '/' {
		urlPathPrefix = "/" + urlPathPrefix
	}
	if urlPathPrefix[len(urlPathPrefix)-1] != '/' {
		urlPathPrefix = urlPathPrefix + "/"
	}
	return urlPathPrefix
}

// extractTargetContentPath return URL path after urlPathPrefixLen as targetContentPath.
// The defaultContentPath will be return if targetContentPath is empty or slash.
func extractTargetContentPath(r *http.Request, urlPathPrefixLen int, defaultContentPath string) (targetContentPath string) {
	if len(r.URL.Path) > urlPathPrefixLen {
		targetContentPath = r.URL.Path[urlPathPrefixLen:]
	}
	if (targetContentPath == "/") || (targetContentPath == "") {
		targetContentPath = defaultContentPath
	}
	return
}
