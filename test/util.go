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

	"github.com/drdgvhbh/gitserver/internal"
	"github.com/drdgvhbh/gitserver/internal/response"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type simpleTestSuite struct {
	suite.Suite
	rootHandler http.Handler
	testServer  *httptest.Server
	basePath    string
	currentDir  string
}

func (suite *simpleTestSuite) SetupTest() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	repoLocation := fmt.Sprintf("%s/simple-git-repo", path.Dir(filename))
	vm, err := createVirtualMachine()
	suite.NoError(err)
	err = vm.cloneRepositoryFrom(repoLocation)
	suite.NoError(err)

	rootHandler := internal.NewRootHandler(osfs.New(""))
	testServer := httptest.NewServer(rootHandler)

	basePath := fmt.Sprintf("%s", strings.Replace(repoLocation, "/", "|", -1))

	suite.rootHandler = rootHandler
	suite.testServer = testServer
	suite.basePath = basePath
	suite.currentDir = path.Dir(repoLocation)
}

type virtualMachine struct {
	fileSystem billy.Filesystem
	storage    *filesystem.Storage
}

func createVirtualMachine() (*virtualMachine, error) {
	directory, err := ioutil.TempDir("", "git-vm")
	if err != nil {
		return nil, err
	}

	fs := osfs.New(directory)
	if err = fs.MkdirAll(".git", 775); err != nil {
		return nil, err
	}

	dotgit, err := fs.Chroot(".git")
	if err != nil {
		return nil, err
	}

	storage := filesystem.NewStorageWithOptions(
		dotgit,
		cache.NewObjectLRUDefault(),
		filesystem.Options{KeepDescriptors: true})

	return &virtualMachine{
		fileSystem: fs,
		storage:    storage,
	}, nil
}

func (vm *virtualMachine) cloneRepositoryFrom(location string) error {
	repoURL, err := url.Parse(location)
	if err != nil {
		return err
	}

	_, err = git.Clone(vm.storage, vm.fileSystem, &git.CloneOptions{
		URL: repoURL.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

type JSONData = []byte

type Suite interface {
	NoError(error, ...interface{}) bool
}

func executeRequest(
	testServer *httptest.Server,
	req *http.Request,
	suite Suite,
) JSONData {
	req.Header.Add("Authorization", "e8b8dc29-d1d9-495d-b509-4dde3701018b")
	resp, err := testServer.Client().Do(req)
	suite.NoError(err)

	body, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)

	var responseBody *response.Definition
	err = json.Unmarshal(body, &responseBody)
	suite.NoError(err)

	data, err := json.Marshal(responseBody.Data)
	suite.NoError(err)

	return data
}
