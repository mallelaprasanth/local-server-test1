package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/local-server-test1/database"
)

func InitDB() {
	database.TestInit()
}

func PerformRequestWithoutAuth(r *gin.Engine, method string, path string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))

	if err != nil {
		panic("Failed to make new HTTP Request.")
	}

	w := httptest.NewRecorder()

	if body != nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	r.ServeHTTP(w, req)
	return w
}
