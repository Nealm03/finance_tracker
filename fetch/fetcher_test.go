package fetch_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Nealm03/finance_tracker/fetch"
	"github.com/stretchr/testify/suite"
)

type FetchTestSuite struct {
	suite.Suite
}

func (fts *FetchTestSuite) Test_FetchFails_WhenAnInvalidUrlIsProvided() {
	url := url.URL{}

	_, err := fetch.JsonFetch[any](context.Background(), url, http.MethodGet)
	fts.ErrorContains(err, "invalid url")
}

func (fts *FetchTestSuite) Test_FetchFails_WhenAnInvalidMethodIsProvided() {
	url, _ := url.Parse("https://google.com")

	_, err := fetch.JsonFetch[any](context.Background(), *url, "")
	fts.ErrorContains(err, "invalid method")
}

func (fts *FetchTestSuite) Test_FetchFails_WhenContextDeadlineExceeded() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(-time.Second*5))
	httptest.NewRecorder()
	_, err := fetch.JsonFetch[any](ctx, *url, http.MethodGet)
	fts.ErrorContains(err, "context deadline exceeded")
}

func (fts *FetchTestSuite) Test_FetchFails_WhenHostReturnsNon200Code() {
	type testcase struct {
		name        string
		code        int
		body        string
		expectedErr string
	}

	cases := []testcase{
		{
			name:        "Bad gateway",
			code:        http.StatusBadGateway,
			body:        ``,
			expectedErr: "server returned a non-successful response code",
		},
		{
			name:        "Bad gateway",
			code:        http.StatusBadRequest,
			body:        ``,
			expectedErr: "server returned a non-successful response code",
		},
		{
			name:        "Bad gateway",
			code:        http.StatusForbidden,
			body:        ``,
			expectedErr: "server returned a non-successful response code",
		},
		{
			name:        "Bad gateway",
			code:        http.StatusUnprocessableEntity,
			body:        ``,
			expectedErr: "server returned a non-successful response code",
		},
	}

	for _, tc := range cases {
		fts.Run(tc.name, func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.code)
				w.Write([]byte(tc.body))
			}))

			defer ts.Close()

			url, _ := url.Parse(ts.URL)

			httptest.NewRecorder()
			_, err := fetch.JsonFetch[any](context.Background(), *url, http.MethodGet)
			fts.ErrorContains(err, tc.expectedErr)
		})
	}
}

func (fts *FetchTestSuite) Test_FetchFails_WhenHostReturnsMalformedJson() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo"))
	}))

	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	httptest.NewRecorder()
	_, err := fetch.JsonFetch[any](context.Background(), *url, http.MethodGet)
	fts.ErrorContains(err, "failed to parse response")
}

func (fts *FetchTestSuite) Test_FetchFails_WhenInvalidTypeProvided() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok": "true"}`))
	}))

	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	httptest.NewRecorder()
	_, err := fetch.JsonFetch[int](context.Background(), *url, http.MethodGet)
	fts.ErrorContains(err, "failed to parse response")
}

func (fts *FetchTestSuite) Test_Returns_UnMarshalsStruct_WhenServerRespondsWithSuccess() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok": true}`))
	}))

	defer ts.Close()
	type TestResp struct {
		Ok bool `json:"ok"`
	}
	url, _ := url.Parse(ts.URL)

	httptest.NewRecorder()
	t, err := fetch.JsonFetch[TestResp](context.Background(), *url, http.MethodGet)

	fts.NoError(err)

	fts.True(t.Ok)
}

func TestFetchSuite(t *testing.T) {
	suite.Run(t, new(FetchTestSuite))
}
