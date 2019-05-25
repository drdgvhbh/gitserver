package internal_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/drdgvhbh/gitserver/internal/mock"

	"github.com/drdgvhbh/gitserver/internal"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func TestContentTypeMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	}

	assert := assert.New(t)

	handler := internal.ContentTypeMiddleware(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)

	assert.EqualValues(res.Header.Get("Content-Type"), "application/json")
}

func TestRequestIDContextMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id")
		fmt.Fprintf(w, "%s", id)
	}

	assert := assert.New(t)

	handler := internal.RequestIDContextMiddleware(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	id := string(body)

	assert.Regexp(uuidRegex, id)
}

func TestRequestMethodContextMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		method := r.Context().Value("method")
		fmt.Fprintf(w, "%s", method)
	}

	assert := assert.New(t)

	const apiVersion = "v1"
	const route = "testing"

	handler := internal.RequestMethodContextMiddleware(
		http.HandlerFunc(mockHandler))
	router := mux.NewRouter().PathPrefix(
		fmt.Sprintf("/%s", apiVersion)).Subrouter()
	router.HandleFunc(fmt.Sprintf("/%s", route), handler.ServeHTTP)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s/%s", ts.URL, apiVersion, route))
	assert.NoError(err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	method := string(body)

	assert.EqualValues(fmt.Sprintf("%s.get", route), method)
}

func TestResponseWriterMiddleware(t *testing.T) {
	testMessage := "testing"
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testMessage)
	}

	writer := httptest.NewRecorder()
	writerFactory := func(http.ResponseWriter, *http.Request) http.ResponseWriter {
		return writer
	}

	assert := assert.New(t)

	handler := internal.NewResponseWriterMiddleware(
		writerFactory)(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	assert.NoError(err)

	resp := writer.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)

	message := string(body)

	assert.EqualValues(testMessage, message)
}

func TestRepositoryDirectoryVariableSanitizer(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	}

	assert := assert.New(t)

	handler := internal.RepositoryDirectoryVariableSanitizer(
		http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req = mux.SetURLVars(req, map[string]string{"directory": "~|home|gitserver"})
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal("~/home/gitserver", mux.Vars(req)["directory"])
}

type NewOpenRepositoryMiddlewareTestSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
	handler  http.Handler
}

func (suite *NewOpenRepositoryMiddlewareTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	})

}

func (suite *NewOpenRepositoryMiddlewareTestSuite) TestReportsRepositoryNotFound() {
	path := "/"
	reader := new(mock.Reader)
	reader.On("Open", path).Return(
		nil, errors.New("repository does not exist"))

	handler := internal.NewOpenRepositoryMiddleware(reader)(suite.handler)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req = mux.SetURLVars(req, map[string]string{"directory": path})
	res := suite.recorder

	handler.ServeHTTP(res, req)

	resp := suite.recorder.Result()
	body, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)

	suite.Assert().Equal(http.StatusNotFound, resp.StatusCode)

	suite.Assert().JSONEq(
		`{ "data": null, "error": { "error": "repository does not exist" }}`,
		string(body))

	reader.AssertExpectations(suite.T())
}
func TestNewOpenRepositoryMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(NewOpenRepositoryMiddlewareTestSuite))
}
