package util

import (
	"os"

	"github.com/golang/glog"
)

// ErrExit prints error log then exit by non-zero.
func ErrExit(err error) {
	glog.Errorln(err)
	os.Exit(1)
}
