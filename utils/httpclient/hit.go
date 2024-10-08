package httpclient

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/smartpet/websocket/metrics"
	"github.com/smartpet/websocket/utils"
)

func (c CustomHttpClient) RequestWithRetries(ctx context.Context, url, method string, headers map[string]string, body []byte,
	retryCount int, retryWaitTime time.Duration, retryMaxWaitTime time.Duration) (*http.Response, error) {

	urlWithoutQuery := utils.GetUrlWithoutQueryParams(url)
	timer := metrics.GetExternalHTTPRequestTimer(urlWithoutQuery)
	defer timer.ObserveDuration()

	var resp *http.Response
	var err error
	resp, err = c.doWithRetries(ctx, method, url, headers, body, retryCount, retryWaitTime, retryMaxWaitTime)
	if resp != nil {
		metrics.SetExternalHTTPRequestCounter(strconv.Itoa(resp.StatusCode), urlWithoutQuery, method)
	}
	return resp, err
}
