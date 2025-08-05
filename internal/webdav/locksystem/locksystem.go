package locksystem

import (
	"golang.org/x/net/webdav"
	"path/filepath"
	"time"
)

var _ webdav.LockSystem = &lockSystem{}

type lockSystem struct {
	webdav.LockSystem
	directory string
}

func (l *lockSystem) Create(now time.Time, details webdav.LockDetails) (token string, err error) {
	details.Root = filepath.Join(l.directory, details.Root)
	return l.LockSystem.Create(now, details)
}

func (l *lockSystem) Confirm(now time.Time, name0, name1 string, conditions ...webdav.Condition) (release func(), err error) {
	if name0 != "" {
		name0 = filepath.Join(l.directory, name0)
	}

	if name1 != "" {
		name1 = filepath.Join(l.directory, name1)
	}

	return l.LockSystem.Confirm(now, name0, name1, conditions...)
}
