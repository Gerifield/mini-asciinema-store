package ministore

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadHandlerInvalidForm(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	s := Server{uploadBucket: bucket}
	srv := httptest.NewServer(s.Routes())
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/asciicasts", strings.NewReader("test data"))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	checkResponse(t, resp, http.StatusBadRequest, "file read failed\n")
}

func TestUploadHandlerSuccess(t *testing.T) {
	bucket := createMemBucketWithFile(t)
	defer bucket.Close()

	s := Server{uploadBucket: bucket}
	srv := httptest.NewServer(s.Routes())
	defer srv.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("asciicast", "asciicast")
	fw.Write([]byte("hi!"))
	require.NoError(t, w.Close())

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/asciicasts", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	rb, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	parts := strings.Split(string(rb), "/")
	require.Equal(t, 3, len(parts))
	assert.Equal(t, "a", parts[1])
	_, err = uuid.Parse(parts[2])
	require.NoError(t, err)
}
