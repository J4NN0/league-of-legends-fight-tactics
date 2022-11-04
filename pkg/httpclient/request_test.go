package httpclient

import (
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/httpclient/httpclienttest"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStruct struct {
	Field1 string
	Field2 string
}

func TestGetSuccess(t *testing.T) {
	path := "some/path"
	expected := MockStruct{Field1: "Field1", Field2: "Field2"}
	dest := MockStruct{}

	client := httpclienttest.NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), path)
		return httpclienttest.Response(expected, 200)
	})

	_ = Get(client, path, &dest)

	assert.Equal(t, expected, dest)
}

func TestGetFail_Not200(t *testing.T) {
	client := httpclienttest.NewTestClient(func(req *http.Request) *http.Response {
		return httpclienttest.Response([]string{}, 400)
	})
	dest := struct{}{}
	err := Get(client, "", &dest)

	assert.Contains(t, err.Error(), "HTTP status not OK")
}

func TestGetFail_UnmarshalWrongStruct(t *testing.T) {
	path := "some/path"
	dest := MockStruct{}

	client := httpclienttest.NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), path)
		return httpclienttest.Response([]byte("xxx"), 200)
	})

	err := Get(client, path, &dest)

	assert.NotNil(t, err)
}
