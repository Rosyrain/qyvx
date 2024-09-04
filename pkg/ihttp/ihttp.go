package ihttp

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

var client = http.Client{Timeout: 10 * time.Second}

func Request(method, url, contentType, token string, data io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		zap.L().Warn("http status code", zap.Any("url", url), zap.Any("code", resp.StatusCode))
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
