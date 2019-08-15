package httpservefile

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
