package util

import (
	"fmt"
	"strings"
)

// ParseRepo splits 'repo' to 'repoOwner' and 'repoName'
// e.g., repo 'wangyoucao577/assets-uploader' => repoOwner 'wangyoucao577' and repoName 'assets-uploader'
func ParseRepo(repo string) (repoOwner string, repoName string, err error) {
	s := strings.Split(repo, "/")
	if len(s) != 2 {
		err = fmt.Errorf("repo has to be 'owner_name/repo_name' format, but got %s", repo)
		return
	}
	repoOwner, repoName = s[0], s[1]
	return
}
