package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func HttpRequester(
	client *http.Client,
	url string,
	headers map[string]string,
	method string,
	data interface{}) (*http.Response, error) {
	var reqBody []byte
	if data != nil {
		var err error
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	return res, nil
}
