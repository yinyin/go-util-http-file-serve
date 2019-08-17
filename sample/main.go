package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	httpservefile "github.com/yinyin/go-util-http-serve-file"
)

func parseCommandParam() (httpAddr, contentFolderPath, contentZipPath, contentStorePrefix, targetRelPath string) {
	flag.StringVar(&httpAddr, "listen", ":8080", "port and address to listen on")
	flag.StringVar(&contentFolderPath, "folder", "", "path to content folder")
	flag.StringVar(&contentZipPath, "zip", "", "path to content zip archive")
	flag.StringVar(&contentStorePrefix, "prefix", "", "content prefix path in content store")
	flag.StringVar(&targetRelPath, "target", "", "relative path to target file")
	flag.Parse()
	return
}

func setupContentServer(contentFolderPath, contentZipPath, contentStorePrefix string) httpservefile.HTTPFileServer {
	switch {
	case contentFolderPath != "":
		srv, err := httpservefile.NewServeFileSystemWithPrefix(contentURLPathPrefix, contentFolderPath)
		if nil != err {
			log.Fatalf("ERROR: failed on setting up content server with FileSystem [%s]: %v", contentFolderPath, err)
			return nil
		}
		log.Printf("INFO: setup content server with FileSystem: [%s]", contentFolderPath)
		return srv
	case contentZipPath != "":
		srv, err := httpservefile.NewServeZipArchiveWithPrefix(contentURLPathPrefix, contentZipPath, contentStorePrefix, "index.html")
		if nil != err {
			log.Fatalf("ERROR: failed on setting up content server with ZipArchive [%s::%s]: %v", contentZipPath, contentStorePrefix, err)
			return nil
		}
		log.Printf("INFO: setup content server with ZipArchive: [%s::%s]", contentZipPath, contentStorePrefix)
		return srv
	default:
		log.Printf("WARN: empty content store")
	}
	return nil
}

func main() {
	httpAddr, contentFolderPath, contentZipPath, contentStorePrefix, targetRelPath := parseCommandParam()
	log.Printf("INFO: listen on address: [%s]", httpAddr)
	h := &sampleHandler{
		contentServer:        setupContentServer(contentFolderPath, contentZipPath, contentStorePrefix),
		targetContentRelPath: targetRelPath,
	}
	defer h.Close()
	log.Printf("INFO: relative path to target content: [%s]", targetRelPath)
	s := &http.Server{
		Addr:         httpAddr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
