package main

import (
	"flag"
	"fmt"
)

type uploaderFlags struct {
	repo   string // e.g., wangyoucao577/assets-uploader
	tag    string
	token  string
	dryrun bool
}

func (u *uploaderFlags) validate() error {
	if len(u.repo) == 0 {
		return fmt.Errorf("repo is mandatory but not set")
	}
	if u.tag == "" {
		return fmt.Errorf("tag or releasename is mandatory but not set")
	}
	if len(u.token) == 0 {
		return fmt.Errorf("github token is mandatory but not set")
	}

	return nil
}

var flags uploaderFlags

func init() {
	flag.StringVar(&flags.repo, "repo", "", "Github repo, e.g., 'wangyoucao577/assets-uploader'.")
	flag.StringVar(&flags.tag, "tag", "", "Git tag to identify a Github Release in repo.")
	flag.StringVar(&flags.token, "token", "", "Github token to make changes.")
	flag.BoolVar(&flags.dryrun, "dryrun", false, "Dry run")
}
