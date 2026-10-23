package tests

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

func sendAndGetContent(t *testing.T, url string, method string, body []byte) (int, string, error) {
	t.Helper()

	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return -1, "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	respBody, _ := io.ReadAll(resp.Body)

	return statusCode, string(respBody), nil
}
