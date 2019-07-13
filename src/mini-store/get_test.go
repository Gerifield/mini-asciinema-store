package ministore

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/memblob"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testEmptyUUID = "f57de750-cb2d-470d-9e45-4b802178e1e6"
	testUUID      = "e26e4c57-7bb1-4969-947b-4ee615758d70"
)

func TestInputParams(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	s := Server{uploadBucket: bucket}
	srv := httptest.NewServer(s.Routes())
	defer srv.Close()

	testTable := []struct {
		URL        string
		StatusCode int
		Message    string
	}{
		{srv.URL + "/a/not-uuid", http.StatusBadRequest, "invalid id\n"},
		{srv.URL + "/a/" + testEmptyUUID, http.StatusNotFound, "file not found\n"},
		{srv.URL + "/a/" + testUUID, http.StatusOK, "hi!"},
	}

	for _, tt := range testTable {
		req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		checkResponse(t, resp, tt.StatusCode, tt.Message)
	}
}

func createMemBucketWithFile(t *testing.T) *blob.Bucket {
	bucket, err := blob.OpenBucket(context.Background(), "mem://")
	require.NoError(t, err)

	err = bucket.WriteAll(context.Background(), testUUID, []byte("hi!"), nil)
	require.NoError(t, err)
	return bucket
}

func checkResponse(t *testing.T, resp *http.Response, statusCode int, message string) {
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, message, string(b))
}
