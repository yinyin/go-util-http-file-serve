package httpservefile

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeFileSystem is an implementation of HTTPFileServer with file system as
// content backend.
type ServeFileSystem struct {
	urlPathPrefixLen  int
	contentFolderPath string
}

// NewServeFileSystemWithPrefixLength create an instance of ServeFileSystem
// with urlPathPrefixLen and contentFolderPath.
func NewServeFileSystemWithPrefixLength(urlPathPrefixLen int, contentFolderPath string) (s *ServeFileSystem, err error) {
	if contentFolderPath == "" {
		return nil, ErrEmptyContentStoragePath
	}
	if contentFolderPath, err = filepath.Abs(contentFolderPath); nil != err {
		return
	}
	if urlPathPrefixLen < 1 {
		urlPathPrefixLen = 1
	}
	return &ServeFileSystem{
		urlPathPrefixLen:  urlPathPrefixLen,
		contentFolderPath: contentFolderPath,
	}, nil
}

// NewServeFileSystemWithPrefix create an instance of ServeFileSystem with
// urlPathPrefix and contentFolderPath.
//
// ** CAUTION **:
// Prefix of URL path will NOT be check. Make sure such check is done at routing logic.
func NewServeFileSystemWithPrefix(urlPathPrefix, contentFolderPath string) (s *ServeFileSystem, err error) {
	urlPathPrefix = sanitizeURLPathPrefix(urlPathPrefix)
	return NewServeFileSystemWithPrefixLength(len(urlPathPrefix), contentFolderPath)
}

// NewServeFileSystem create an instance of ServeFileSystem with contentFolderPath.
func NewServeFileSystem(contentFolderPath string) (s *ServeFileSystem, err error) {
	return NewServeFileSystemWithPrefixLength(1, contentFolderPath)
}

func (s *ServeFileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request, defaultFileName, targetFileName string) {
	if targetFileName == "" {
		if len(r.URL.Path) > s.urlPathPrefixLen {
			targetFileName = r.URL.Path[s.urlPathPrefixLen:]
		}
		if (targetFileName == "/") == (targetFileName == "") {
			targetFileName = defaultFileName
		}
	}
	targetFilePath := filepath.Join(s.contentFolderPath, targetFileName)
	if !strings.HasPrefix(targetFilePath, s.contentFolderPath) {
		http.NotFound(w, r)
		return
	}
	fileinfo, err := os.Stat(targetFilePath)
	if nil != err {
		http.Error(w, "internal error (file-system)", http.StatusInternalServerError)
		log.Printf("WARN: failed on stat file [%s]: %v", targetFilePath, err)
		return
	}
	if fileinfo.IsDir() {
		http.NotFound(w, r)
		return
	}
	fp, err := os.Open(targetFilePath)
	if nil != err {
		http.Error(w, "internal error (file-system)", http.StatusInternalServerError)
		log.Printf("WARN: failed on open file [%s]: %v", targetFilePath, err)
		return
	}
	defer fp.Close()
	http.ServeContent(w, r, targetFilePath, fileinfo.ModTime(), fp)
}

// Close free used resources.
func (s *ServeFileSystem) Close() (err error) {
	return
}
