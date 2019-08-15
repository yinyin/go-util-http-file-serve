package httpservefile

import (
	"net/http"
	"net/url"

	utilhttphandlers "github.com/yinyin/go-util-http-handlers"
)

// ServeZipArchive is an implementation of HTTPFileServer with a zip file as
// content backend.
type ServeZipArchive struct {
	urlPathPrefixLen int
	zipContentServer *utilhttphandlers.ZipArchiveContentServer
}

// NewServeZipArchiveWithPrefixLength create an instance of ServeZipArchive
// with urlPathPrefixLen, zipFilePath, zipPathPrefix and zipDefaultContentPath.
func NewServeZipArchiveWithPrefixLength(urlPathPrefixLen int, zipFilePath, zipPathPrefix, zipDefaultContentPath string) (s *ServeZipArchive, err error) {
	if zipFilePath == "" {
		return nil, ErrEmptyContentStoragePath
	}
	zipContentServer, err := utilhttphandlers.NewZipArchiveContentServer(zipFilePath, zipPathPrefix, zipDefaultContentPath)
	if nil != err {
		return
	}
	if urlPathPrefixLen < 1 {
		urlPathPrefixLen = 1
	}
	return &ServeZipArchive{
		urlPathPrefixLen: urlPathPrefixLen,
		zipContentServer: zipContentServer,
	}, nil
}

// NewServeZipArchiveWithPrefix create an instance of ServeZipArchive with
// urlPathPrefix, zipFilePath, zipPathPrefix and zipDefaultContentPath.
//
// ** CAUTION **:
// Prefix of URL path will NOT be check. Make sure such check is done at routing logic.
func NewServeZipArchiveWithPrefix(urlPathPrefix, zipFilePath, zipPathPrefix, zipDefaultContentPath string) (s *ServeZipArchive, err error) {
	urlPathPrefix = sanitizeURLPathPrefix(urlPathPrefix)
	return NewServeZipArchiveWithPrefixLength(len(urlPathPrefix), zipFilePath, zipPathPrefix, zipDefaultContentPath)
}

// NewServeZipArchive create an instance of ServeFileSystem with zipFilePath, zipPathPrefix and zipDefaultContentPath.
func NewServeZipArchive(zipFilePath, zipPathPrefix, zipDefaultContentPath string) (s *ServeZipArchive, err error) {
	return NewServeZipArchiveWithPrefixLength(1, zipFilePath, zipPathPrefix, zipDefaultContentPath)
}

func (s *ServeZipArchive) ServeHTTP(w http.ResponseWriter, r *http.Request, defaultFileName, targetFileName string) {
	if nil == s.zipContentServer {
		http.Error(w, "archive file closed", http.StatusBadGateway)
		return
	}
	if targetFileName == "" {
		if len(r.URL.Path) > s.urlPathPrefixLen {
			targetFileName = r.URL.Path[s.urlPathPrefixLen:]
		}
		if (targetFileName == "/") == (targetFileName == "") {
			targetFileName = defaultFileName
		}
	}
	if targetFileName == "" {
		http.NotFound(w, r)
		return
	}
	if targetFileName[0] != '/' {
		targetFileName = "/" + targetFileName
	}
	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = targetFileName
	s.zipContentServer.ServeHTTP(w, r2)
}

// Close free used resources.
func (s *ServeZipArchive) Close() (err error) {
	if nil == s.zipContentServer {
		return
	}
	err = s.zipContentServer.Close()
	s.zipContentServer = nil
	return
}
