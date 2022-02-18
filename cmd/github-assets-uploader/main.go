package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wangyoucao577/assets-uploader/util/appversion"

	"github.com/golang/glog"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func errExit(err error) {
	glog.Errorln(err)
	os.Exit(1)
}

func parseRepo(repo string) (repoOwner string, repoName string, err error) {
	s := strings.Split(repo, "/")
	if len(s) != 2 {
		err = fmt.Errorf("repo has to be 'owner_name/repo_name' format, but got %s", repo)
		return
	}
	repoOwner, repoName = s[0], s[1]
	return
}

func main() {
	flag.Parse()
	appversion.PrintExit()
	defer glog.Flush()

	if err := flags.validate(); err != nil {
		errExit(err)
	}
	repoOwner, repoName, err := parseRepo(flags.repo)
	if err != nil {
		errExit(err)
	}

	retry := flags.retry
	for {
		retry--

		err = uploadAsset(repoOwner, repoName, flags.tag, flags.file, flags.mediaType, flags.token, flags.overwrite)
		if err != nil {
			if retry == 0 {
				errExit(err)
			} else {
				glog.Warningf("Upload asset error, will retry soon: %v\n", err)
				time.Sleep(3 * time.Second) // retry after 3 seconds
				continue
			}
		}

		break // break when succeed
	}
}

func uploadAsset(repoOwner, repoName, tag, assetPath, mediaType, token string, overwrite bool) error {

	// read-write client
	rwContext := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(rwContext, ts)
	client := github.NewClient(tc)

	// get release by tag
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), repoOwner, repoName, tag)
	if err != nil {
		return err
	}

	assetName := filepath.Base(assetPath)
	if overwrite { // remove old one if it's exist already
		assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), repoOwner, repoName, release.GetID(), nil)
		if err != nil {
			return err
		}
		for _, asset := range assets {
			if asset.GetName() == assetName {

				// found exist one, delete it
				if _, err := client.Repositories.DeleteReleaseAsset(rwContext, repoOwner, repoName, asset.GetID()); err != nil {
					return err
				}
				glog.Infof("Deleted old asset, id %d, name '%s', url '%s'\n", asset.GetID(), asset.GetName(), asset.GetBrowserDownloadURL())
				break
			}
		}
	}

	// open file for uploading
	f, err := os.Open(assetPath) // For read access.
	if err != nil {
		return err
	}
	defer f.Close()

	// upload
	releaseAsset, _, err := client.Repositories.UploadReleaseAsset(rwContext, repoOwner, repoName, release.GetID(), &github.UploadOptions{
		Name:      assetName,
		Label:     "",
		MediaType: mediaType,
	}, f)
	if err != nil {
		return err
	}
	glog.Infof("Upload asset succeed, id %d, name '%s', url: '%s'\n", releaseAsset.GetID(), releaseAsset.GetName(), releaseAsset.GetBrowserDownloadURL())
	return nil
}
