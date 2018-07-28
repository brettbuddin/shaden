package runtime

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettbuddin/shaden/errors"
	"github.com/stretchr/testify/require"
)

func TestHandler_GoodEval(t *testing.T) {
	var (
		mux = http.NewServeMux()
		e   evaler
	)

	AddHandler(mux, &e)
	s := httptest.NewServer(mux)
	defer s.Close()

	client := s.Client()
	url := s.URL + "/eval"

	const content = "content"

	buf := bytes.NewBufferString(content)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, "OK", string(body))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, content, string(e.content))
}

func TestHandler_BadEval(t *testing.T) {
	var (
		mux = http.NewServeMux()
		e   = evaler{
			err: errors.New("bad"),
		}
	)

	AddHandler(mux, &e)
	s := httptest.NewServer(mux)
	defer s.Close()

	client := s.Client()
	url := s.URL + "/eval"

	const content = "content"

	buf := bytes.NewBufferString(content)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, "bad", string(body))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, content, string(e.content))
}

func TestHandler_UnimplementedMethod(t *testing.T) {
	var (
		mux = http.NewServeMux()
		e   evaler
	)

	AddHandler(mux, &e)
	s := httptest.NewServer(mux)
	defer s.Close()

	client := s.Client()
	url := s.URL + "/eval"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.StatusCode)
}

type evaler struct {
	content []byte
	val     interface{}
	err     error
}

func (e *evaler) Eval(b []byte) (interface{}, error) {
	e.content = b
	return e.val, e.err
}
