package fixtures

import (
	"bytes"
	"net/http"
)

func CreateFakeRequest(method, path, auth string, body []byte) *http.Request {
	var request *http.Request

	if body != nil {
		request, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		request, _ = http.NewRequest(method, path, nil)
	}

	if auth != "" {
		request.Header.Add("Authorization", auth)
	}

	return request
}
