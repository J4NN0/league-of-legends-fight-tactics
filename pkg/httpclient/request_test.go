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

func TestGetFailNot200(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return mockResponse([]string{}, 400)
	})
	dest := struct{}{}
	err := Get(client, "", &dest)

	assert.Contains(t, err.Error(), "HTTP status not OK")
}

func TestExtractItemSuccess(t *testing.T) {
	expected := MockStruct{Field1: "Field1", Field2: "Field2"}
	dest := MockStruct{}

	expectedJsn, _ := json.Marshal(expected)
	_ = unmarshalToInterface(expectedJsn, &dest)

	assert.Equal(t, expected, dest)
}

func TestExtractItemWrongStruct(t *testing.T) {
	notMyMock := struct {
		SomeField string
	}{
		SomeField: "WhatAField",
	}
	myMock := MockStruct{}

	notMyMockJsn, _ := json.Marshal(notMyMock)
	_ = unmarshalToInterface(notMyMockJsn, &myMock)

	assert.Equal(t, MockStruct{}, myMock)
}

func TestExtractItemFail(t *testing.T) {
	dest := ""
	err := unmarshalToInterface([]byte("xxx"), &dest)

	assert.NotNil(t, err)
}
