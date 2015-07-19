package testutil

import (
	"net/http"
	"net/http/httptest"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type mockHandler struct {
	h HandlerFunc
}

func (m mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.h(w, r)
}

func NewMockServer(h HandlerFunc) *httptest.Server {
	return httptest.NewServer(mockHandler{h})
}
