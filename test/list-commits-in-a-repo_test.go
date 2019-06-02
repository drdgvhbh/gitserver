package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ListCommitsInARepoTestSuite struct {
	simpleTestSuite
}

func (suite *ListCommitsInARepoTestSuite) TestListCommits() {
	testServer := suite.testServer
	basePath := suite.basePath
	currentDir := suite.currentDir

	defer testServer.Close()

	reqURL, err := url.Parse(
		fmt.Sprintf("%s/v1/repositories/%s/commits", testServer.URL, basePath))
	suite.NoError(err)

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	suite.NoError(err)

	data := executeRequest(testServer, req, suite)

	testDataFilePath := fmt.Sprintf(
		"%s/test-responses/list-commits-simple.json",
		currentDir)
	testData, err := ioutil.ReadFile(testDataFilePath)
	suite.NoError(err)

	suite.Assert().JSONEq(
		`[
			{
				"hash": "be50985852e7aadc4392fb4809f3f9e265a92694",
				"summary": "Merge branch 'branch'",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-27T23:11:34-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-27T23:11:34-04:00"
				}
			},
			{
				"hash": "a7170f7640bb9b9960fe8a20b4454f71f98c423d",
				"summary": "Branch",
				"author": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-27T23:08:27-04:00"
				},
				"committer": {
					"name": "Ryan Lee",
					"email": "drdgvhbh@gmail.com",
					"timestamp": "2019-05-27T23:08:27-04:00"
				}
			},
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
