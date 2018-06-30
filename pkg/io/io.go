package io

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func RedirectUpdate(updater rsrc.Updater) *updateRedirect {
	return &updateRedirect{updater: updater}
}

func (ur updateRedirect) Read(loc rsrc.Locator) (data []byte, err error) {
	return ur.updater.Update(loc)
}
