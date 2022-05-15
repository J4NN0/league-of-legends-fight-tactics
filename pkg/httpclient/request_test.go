package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func mockResponse(obj interface{}, status int) *http.Response {
	jsonMarshal, _ := json.Marshal(obj)
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewReader(jsonMarshal)),
		Header:     make(http.Header),
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

type MockStruct struct {
	Field1 string
	Field2 string
}

func TestGetSuccess(t *testing.T) {
	path := "some/path"
	expected := MockStruct{Field1: "Field1", Field2: "Field2"}
	dest := MockStruct{}

	client := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), path)
		return mockResponse(expected, 200)
	})

	_ = Get(client, path, &dest)

	assert.Equal(t, expected, dest)
}

func TestGetFail_Not200(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return mockResponse([]string{}, 400)
	})
	dest := struct{}{}
	err := Get(client, "", &dest)

	assert.Contains(t, err.Error(), "HTTP status not OK")
}

func TestGetFail_UnmarshalWrongStruct(t *testing.T) {
	path := "some/path"
	dest := MockStruct{}

	client := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), path)
		return mockResponse([]byte("xxx"), 200)
	})

	err := Get(client, path, &dest)

	assert.NotNil(t, err)
}
