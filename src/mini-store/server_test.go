package ministore

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob"
)

func TestNew(t *testing.T) {
	bucket := &blob.Bucket{}

	s := New("testURL", bucket, "testAuth")

	assert.Equal(t, "testURL", s.baseURL)
	assert.Equal(t, bucket, s.uploadBucket)
	assert.Equal(t, "testAuth", s.authFile)
}

func TestAuthInvalidField(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	srv := httptest.NewServer(authMiddleware([]string{"token1"})(http.NewServeMux()))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost,srv.URL, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	checkResponse(t, resp, http.StatusUnauthorized, "invalid Authorization field\n")
}

func TestAuthInvalidToken(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	srv := httptest.NewServer(authMiddleware([]string{"token1"})(http.NewServeMux()))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost,srv.URL, nil)
	require.NoError(t, err)
	req.SetBasicAuth("test", "tests")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	checkResponse(t, resp, http.StatusUnauthorized, "invalid token\n")
}

func TestAuthSuccess(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	srv := httptest.NewServer(authMiddleware([]string{"token1"})(http.NewServeMux()))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost,srv.URL, nil)
	require.NoError(t, err)
	req.SetBasicAuth("test", "token1")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	checkResponse(t, resp, http.StatusNotFound, "404 page not found\n")
}
