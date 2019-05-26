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
				"hash": "625d85387d80a56a26a5c7ff28d84e49afef2635",
				"summary": "This is me adding text to a file",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T17:19:52-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T17:19:52-04:00"
				}
			},
			{
				"hash": "22b86bcafbc105b45d62772546ef1d2b4bef74f1",
				"summary": "this is me renaming a directory",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:18:47-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:18:47-04:00"
				}
			},
			{
				"hash": "a004980430dffe788c44d97584a2818c54bf6e54",
				"summary": "this is me adding a file in a subdirectory",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:18:15-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:18:15-04:00"
				}
			},
			{
				"hash": "b322ecd55ffee1cdfd42a68f49f201129d9ef023",
				"summary": "this is me deleting a file",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:17:35-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:17:35-04:00"
				}
			},
			{
				"hash": "9649993898da5b19e530a2750ac5fb6fea2d56f8",
				"summary": "This is me adding three new files",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:16:38-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:16:38-04:00"
				}
			},
			{
				"hash": "a82386953aa55fe38dd79a3b65742e39298afc2f",
				"summary": "This is the start of the repo",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:15:37-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-25T15:15:37-04:00"
				}
			}
		]`,
		string(data),
	)

}

func TestListCommitsInARepoTestSuite(t *testing.T) {
	suite.Run(t, new(ListCommitsInARepoTestSuite))
}
