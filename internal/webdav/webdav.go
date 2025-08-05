package webdav

import (
	"golang.org/x/net/webdav"
)

func New(prefix string, path string) *webdav.Handler {
	dav := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(path),
		LockSystem: webdav.NewMemLS(),
		Logger:     nil,
	}

	return dav
}
