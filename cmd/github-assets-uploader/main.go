package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangyoucao577/assets-uploader/util/appversion"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func errExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	appversion.PrintExit()

	if err := flags.validate(); err != nil {
		errExit(err)
	}
	repoInfo := strings.Split(flags.repo, "/")
	if len(repoInfo) != 2 {
		errExit(fmt.Errorf("repo has to be 'owner_name/repo_name' format, but got %s", flags.repo))
	}
	repoOwner := repoInfo[0]
	repoName := repoInfo[1]

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), repoOwner, repoName, flags.tag)
	if err != nil {
		errExit(err)
	}

	f, err := os.Open(flags.file) // For read access.
	if err != nil {
		errExit(err)
	}
	defer f.Close()

	uploadContext := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: flags.token})
	tc := oauth2.NewClient(uploadContext, ts)
	uploadClient := github.NewClient(tc)

	releaseAsset, _, err := uploadClient.Repositories.UploadReleaseAsset(context.Background(), repoOwner, repoName, release.GetID(), &github.UploadOptions{
		Name:      filepath.Base(f.Name()),
		Label:     "",
		MediaType: flags.mediaType,
	}, f)
	if err != nil {
		errExit(err)
	}

	fmt.Printf("Upload successed: %s\n", releaseAsset.GetBrowserDownloadURL())
}
