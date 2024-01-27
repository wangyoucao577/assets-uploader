package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/wangyoucao577/assets-uploader/util"
	"github.com/wangyoucao577/assets-uploader/util/appversion"
)

func main() {
	flag.Parse()
	appversion.PrintExit()
	defer glog.Flush()

	if err := flags.validate(); err != nil {
		util.ErrExit(err)
	}
	repoOwner, repoName, err := util.ParseRepo(flags.repo)
	if err != nil {
		util.ErrExit(err)
	}

	err = cleanAssets(repoOwner, repoName, flags.tag, flags.token, flags.dryrun)
	if err != nil {
		util.ErrExit(err)
	}
}

func cleanAssets(repoOwner, repoName, tag, token string, dryrun bool) error {
	// read-write client
	rwContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(rwContext, ts)
	client := github.NewClient(tc)

	var release *github.RepositoryRelease
	var err error
	release, _, err = client.Repositories.GetReleaseByTag(rwContext, repoOwner, repoName, tag)
	if err != nil {
		return err
	}

	var assets []*github.ReleaseAsset
	assets, _, err = client.Repositories.ListReleaseAssets(rwContext, repoOwner, repoName, release.GetID(), nil)
	if err != nil {
		return err
	}
	for _, asset := range assets {

		// found exist one, delete it
		if !dryrun {
			if _, err = client.Repositories.DeleteReleaseAsset(rwContext, repoOwner, repoName, asset.GetID()); err != nil {
				return err
			}
		}

		glog.Infof("Deleted asset, id %d, name '%s', url '%s'\n", asset.GetID(), asset.GetName(), asset.GetBrowserDownloadURL())
	}
	if len(assets) > 0 {
		glog.Infof("Total deleted %d assets\n", len(assets))
	}

	return nil
}
