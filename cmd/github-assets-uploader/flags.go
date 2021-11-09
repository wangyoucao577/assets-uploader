package main

import (
	"flag"
	"fmt"
)

type uploaderFlags struct {
	file      string
	mediaType string
	repo      string // e.g., wangyoucao577/assets-uploader
	tag       string
	token     string
	overwrite bool
	retry     uint
}

func (u *uploaderFlags) validate() error {
	if len(u.repo) == 0 {
		return fmt.Errorf("repo is mandatory but not set")
	}
	if len(u.tag) == 0 {
		return fmt.Errorf("tag is mandatory but not set")
	}
	if len(u.file) == 0 {
		return fmt.Errorf("file is mandatory but not set")
	}
	if len(u.token) == 0 {
		return fmt.Errorf("github token is mandatory but not set")
	}

	return nil
}

var flags uploaderFlags

func init() {
	flag.StringVar(&flags.file, "f", "", "File path to upload.")
	flag.StringVar(&flags.mediaType, "mediatype", "application/gzip", "E.g., 'application/zip'.")
	flag.StringVar(&flags.repo, "repo", "", "Github repo, e.g., 'wangyoucao577/assets-uploader'.")
	flag.StringVar(&flags.tag, "tag", "", "Git tag to identify a Github Release in repo.")
	flag.StringVar(&flags.token, "token", "", "Github token to make changes.")
	flag.BoolVar(&flags.overwrite, "overwrite", false, "Overwrite release asset if it's already exist.")
	flag.UintVar(&flags.retry, "retry", 1, "How many times to retry if error occur.")
}
