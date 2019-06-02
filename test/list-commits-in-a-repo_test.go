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
		string(testData),
		string(data),
	)

}

func TestListCommitsInARepoTestSuite(t *testing.T) {
	suite.Run(t, new(ListCommitsInARepoTestSuite))
}
