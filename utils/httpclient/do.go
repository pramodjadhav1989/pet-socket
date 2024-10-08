package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	log "github.com/smartpet/websocket/utils/logger"
)

type CustomHttpClient struct {
	*http.Client
}

func (c CustomHttpClient) getRequest(method, url string, headers map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	return request, err
}

func (c CustomHttpClient) shouldRetry(attempt, retryCount int, response *http.Response, err error) bool {
	// check with attempt count
	if attempt >= retryCount {
		return false
	}
	// connect error
	if err != nil {
		return true
	}
	// server error
	if response.StatusCode == 0 || response.StatusCode >= 500 {
		return true
	}
	return false
}

func (c CustomHttpClient) getWaitDuration(attempt int, retryWaitTime time.Duration, retryMaxWaitTime time.Duration) time.Duration {
	waitDuration := retryWaitTime * time.Duration(attempt+1)
	if waitDuration > retryMaxWaitTime {
		return retryMaxWaitTime
	}
	return waitDuration
}

func (c CustomHttpClient) doWithRetries(ctx context.Context, method, url string, headers map[string]string, body []byte, retryCount int, retryWaitTime time.Duration,
	retryMaxWaitTime time.Duration) (*http.Response, error) {
	for attempt := 0; attempt <= retryCount; attempt++ {
		if attempt > 0 {
			log.ApplicationInfo(ctx).Msgf("Retrying count: %v for url: %v and with headers: %v", attempt, url, headers)
		}
		req, err := c.getRequest(method, url, headers, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		response, err := c.Do(req)
		if ctx.Err() != nil {
			return nil, err
		}
		needsRetry := c.shouldRetry(attempt, retryCount, response, err)
		if !needsRetry {
			return response, err
		}
		waitDuration := c.getWaitDuration(attempt, retryWaitTime, retryMaxWaitTime)
		time.Sleep(waitDuration)
	}
	return nil, errors.New("unable to process request")
}
