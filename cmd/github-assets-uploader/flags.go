package main

import (
	"flag"
	"fmt"
)

type uploaderFlags struct {
	file        string
	mediaType   string
	overwrite   bool
	releaseName string
	repo        string // e.g., wangyoucao577/assets-uploader
	retry       uint
	tag         string
	token       string
	baseUrl		string
}

func (u *uploaderFlags) validate() error {
	if len(u.repo) == 0 {
		return fmt.Errorf("repo is mandatory but not set")
	}
	if u.tag == "" && u.releaseName == "" {
		return fmt.Errorf("tag or releasename is mandatory but not set")
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
	flag.BoolVar(&flags.overwrite, "overwrite", false, "Overwrite release asset if it's already exist.")
	flag.StringVar(&flags.releaseName, "releasename", "", "Upload asset to an existing named release.")
	flag.StringVar(&flags.repo, "repo", "", "Github repo, e.g., 'wangyoucao577/assets-uploader'.")
	flag.StringVar(&flags.tag, "tag", "", "Git tag to identify a Github Release in repo.")
	flag.UintVar(&flags.retry, "retry", 1, "How many times to retry if error occur.")
	flag.StringVar(&flags.token, "token", "", "Github token to make changes.")
	flag.StringVar(&flags.baseUrl, "baseurl", "", "Github base URL. E.g. http://github.example.com")
}
