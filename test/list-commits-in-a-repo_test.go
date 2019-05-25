package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	"gopkg.in/src-d/go-billy.v4/osfs"

	"gopkg.in/src-d/go-git.v4"

	"github.com/drdgvhbh/gitserver/internal"
	"github.com/stretchr/testify/suite"
)

type ListCommitsInARepoTestSuite struct {
	suite.Suite
	repo *git.Repository
}

func (suite *ListCommitsInARepoTestSuite) SetupTest() {

}

func (suite *ListCommitsInARepoTestSuite) TestDerp() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	repoURL, err := url.Parse(fmt.Sprintf("%s/simple-git-repo", path.Dir(filename)))
	suite.NoError(err)

	directory, err := ioutil.TempDir("", "basictest")
	suite.NoError(err)
	fs := osfs.New(directory)

	err = fs.MkdirAll(".git", 775)
	suite.NoError(err)

	dotgit, err := fs.Chroot(".git")
	suite.NoError(err)

	storage := filesystem.NewStorageWithOptions(
		dotgit,
		cache.NewObjectLRUDefault(),
		filesystem.Options{KeepDescriptors: true})
	_, err = git.Clone(storage, fs, &git.CloneOptions{
		URL: repoURL.String(),
	})
	suite.NoError(err)

	rootHandler := internal.NewRootHandler(osfs.New(""))
	testServer := httptest.NewServer(rootHandler)

	defer testServer.Close()

	repoPath := fmt.Sprintf("%s", strings.Replace(directory, "/", "|", -1))
	reqURL, err := url.Parse(
		fmt.Sprintf("%s/v1/repositories/%s/commits", testServer.URL, repoPath))
	suite.NoError(err)

	resp, err := http.Get(reqURL.String())
	suite.NoError(err)

	body, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)

	var responseBody *internal.Response
	err = json.Unmarshal(body, &responseBody)
	suite.NoError(err)

	data, err := json.Marshal(responseBody.Data)
	suite.NoError(err)

	suite.Assert().JSONEq(
		`[
			{
				"summary": "This is me adding text to a file"
			},
			{
				"summary": "this is me renaming a directory"
			},
			{
				"summary": "this is me adding a file in a subdirectory"
			},
			{
				"summary": "this is me deleting a file"
			},
			{
				"summary": "This is me adding three new files"
			},
			{
				"summary": "This is the start of the repo"
			}
		]`,
		string(data),
	)

}

func TestListCommitsInARepoTestSuite(t *testing.T) {
	suite.Run(t, new(ListCommitsInARepoTestSuite))
}
