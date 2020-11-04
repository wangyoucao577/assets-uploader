package main

import (
	"flag"
	"fmt"
)

type uploaderFlags struct {
	file      string
	mediaType string
	owner     string // e.g., wangyoucao577
	repo      string // e.g., assets-uploader
	tag       string
	token     string
	overwrite bool
}

func (u *uploaderFlags) validate() error {
	if len(u.owner) == 0 {
		return fmt.Errorf("owner is mandatory but not set")
	}
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
		return fmt.Errorf("Github token is mandatory but not set")
	}

	return nil
}

var flags uploaderFlags

func init() {
	flag.StringVar(&flags.file, "f", "", "File path to upload.")
	flag.StringVar(&flags.mediaType, "mediatype", "application/gzip", "E.g., 'application/zip'.")
	flag.StringVar(&flags.owner, "owner", "", "Github owner.")
	flag.StringVar(&flags.repo, "repo", "", "Github repo.")
	flag.StringVar(&flags.tag, "tag", "", "Git tag to identify a Github Release in repo.")
	flag.StringVar(&flags.token, "token", "", "Github token to make changes.")
	flag.BoolVar(&flags.overwrite, "overwrite", false, "Overwrite release asset if it's already exist.")
}
