package requests

import (
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpRequester(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mocking server response
	url := "https://example.com/api"
	headers := map[string]string{
		"Authorization": "Bearer test-token",
		"Content-Type":  "application/json",
	}

	method := "POST"
	data := map[string]string{
		"key": "value",
	}

	responseBody := `{"success": true}`
	httpmock.RegisterResponder(method, url,
		httpmock.NewStringResponder(200, responseBody))

	client := &http.Client{}

	// Invoke HttpRequester
	resp, err := HttpRequester(client, url, headers, method, data)
	if err != nil {
		t.Fatalf("HttpRequester returned an error: %v", err)

	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 status code, recieved: %d", resp.StatusCode)
	}

	// Check the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Response body reading error: %v", err)

	}

	expectedBody := responseBody
	if string(body) != expectedBody {
		t.Errorf("Expected response: %s, got: %s", expectedBody, string(body))
	}

	// Check if the request was made correctly
	info := httpmock.GetCallCountInfo()
	if count := info[method+" "+url]; count != 1 {
		t.Errorf("Expected 1 request, got: %d", count)
	}
}
