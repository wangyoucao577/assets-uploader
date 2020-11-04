package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func errExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if err := flags.validate(); err != nil {
		errExit(err)
	}

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), flags.owner, flags.repo, flags.tag)
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

	releaseAsset, _, err := uploadClient.Repositories.UploadReleaseAsset(context.Background(), flags.owner, flags.repo, release.GetID(), &github.UploadOptions{
		Name:      filepath.Base(f.Name()),
		Label:     "",
		MediaType: flags.mediaType,
	}, f)
	if err != nil {
		errExit(err)
	}

	fmt.Printf("Upload successed: %s\n", releaseAsset.GetBrowserDownloadURL())
}
