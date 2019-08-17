package main

import (
	"net/http"
	"strings"

	httpservefile "github.com/yinyin/go-util-http-serve-file"
)

const (
	contentURLPathPrefix  = "/content/"
	contentTargetFileName = contentURLPathPrefix + "target"
	defaultFileName       = "index.html"
)

type sampleHandler struct {
	contentServer        httpservefile.HTTPFileServer
	targetContentRelPath string
}

func (h *sampleHandler) Close() (err error) {
	if h.contentServer == nil {
		return
	}
	err = h.contentServer.Close()
	h.contentServer = nil
	return
}

func (h *sampleHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == contentTargetFileName {
		h.contentServer.ServeHTTP(w, req, defaultFileName, h.targetContentRelPath)
		return
	} else if strings.HasPrefix(req.URL.Path, contentURLPathPrefix) {
		h.contentServer.ServeHTTP(w, req, defaultFileName, "")
		return
	}
	http.NotFound(w, req)
}
