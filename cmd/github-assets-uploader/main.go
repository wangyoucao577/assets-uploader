package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/wangyoucao577/assets-uploader/util/appversion"
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
	minDuration := 3
	maxDuration := 15
	for {
		retry--

		err = uploadAsset(repoOwner, repoName, flags.tag, flags.file, flags.mediaType, flags.token, flags.baseUrl, flags.overwrite)
		if err != nil {
			if retry == 0 {
				errExit(err)
			} else {
				randomDuration := time.Duration(minDuration + rand.Intn(maxDuration-minDuration))
				retryDuration := time.Second * randomDuration

				glog.Warningf("Upload asset error, will retry in %s: %v\n", retryDuration.String(), err)
				time.Sleep(retryDuration) // retry after 3-15 seconds
				continue
			}
		}

		break // break when succeed
	}
}

func uploadAsset(repoOwner, repoName, tag, assetPath, mediaType, token string, baseUrl string, overwrite bool) error {
	// read-write client
	rwContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(rwContext, ts)

	var client *github.Client
	if baseUrl != "" {
		client, _ = github.NewEnterpriseClient(baseUrl, baseUrl, tc)
	} else {
		client = github.NewClient(tc)
	}
	

	var release *github.RepositoryRelease
	var err error

	// get release by tag or name
	if tag != "" {
		release, _, err = client.Repositories.GetReleaseByTag(rwContext, repoOwner, repoName, tag)
	} else if flags.releaseName != "" {
		release, err = getReleaseByName(rwContext, client, repoOwner, repoName, flags.releaseName)
	}

	if err != nil {
		return err
	}

	assetName := filepath.Base(assetPath)
	if overwrite { // remove old one if it's exist already
		var assets []*github.ReleaseAsset
		assets, _, err = client.Repositories.ListReleaseAssets(rwContext, repoOwner, repoName, release.GetID(), nil)
		if err != nil {
			return err
		}
		for _, asset := range assets {
			if asset.GetName() == assetName {

				// found exist one, delete it
				if _, err = client.Repositories.DeleteReleaseAsset(rwContext, repoOwner, repoName, asset.GetID()); err != nil {
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

// getReleaseByName lists all releases, sort them by created-at and picks the first/newest one with the given name.
func getReleaseByName(ctx context.Context, client *github.Client, repoOwner, repoName, releaseName string) (*github.RepositoryRelease, error) {
	releases, _, err := client.Repositories.ListReleases(ctx, repoOwner, repoName, nil)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].CreatedAt.After(releases[j].CreatedAt.Time)
	})

	for _, release := range releases { // assume they are some kind of sorted, first pick (newest)
		if release.GetName() == releaseName {
			return release, nil
		}
	}
	return nil, fmt.Errorf("no release found with name: %q", releaseName)
}
