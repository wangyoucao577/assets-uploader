package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyoucao577/assets-uploader/util/appversion"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func errExit(err error) {
	fmt.Println(err)
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

	if err := flags.validate(); err != nil {
		errExit(err)
	}
	repoOwner, repoName, err := parseRepo(flags.repo)
	if err != nil {
		errExit(err)
	}

	// read-write client
	rwContext := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: flags.token})
	tc := oauth2.NewClient(rwContext, ts)
	client := github.NewClient(tc)

	// get release by tag
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), repoOwner, repoName, flags.tag)
	if err != nil {
		errExit(err)
	}

	assetName := filepath.Base(flags.file)
	if flags.overwrite { // remove old one if it's exist already
		assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), repoOwner, repoName, release.GetID(), nil)
		if err != nil {
			errExit(err)
		}
		for _, asset := range assets {
			if asset.GetName() == assetName {

				// found exist one, delete it
				if _, err := client.Repositories.DeleteReleaseAsset(rwContext, repoOwner, repoName, asset.GetID()); err != nil {
					errExit(err)
				}
				fmt.Printf("Deleted old asset, id %d, name '%s', url '%s'\n", asset.GetID(), asset.GetName(), asset.GetBrowserDownloadURL())
				break
			}
		}
	}

	// open file for uploading
	f, err := os.Open(flags.file) // For read access.
	if err != nil {
		errExit(err)
	}
	defer f.Close()

	// upload
	releaseAsset, _, err := client.Repositories.UploadReleaseAsset(rwContext, repoOwner, repoName, release.GetID(), &github.UploadOptions{
		Name:      assetName,
		Label:     "",
		MediaType: flags.mediaType,
	}, f)
	if err != nil {
		errExit(err)
	}
	fmt.Printf("Upload asset succeed, id %d, name '%s', url: '%s'\n", releaseAsset.GetID(), releaseAsset.GetName(), releaseAsset.GetBrowserDownloadURL())
}
