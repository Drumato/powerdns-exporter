package powerdns

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	powerDNSAPIKeyHeader = "X-API-Key"
)

type DefaultClient struct {
	logger  *slog.Logger
	baseURL string
	client  http.Client
	apiKey  string
}

func (d *DefaultClient) GetServers(ctx context.Context) ([]Server, error) {
	url, err := url.JoinPath(d.baseURL, "api", "v1", "servers")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req, err := d.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	servers, err := doDefaultClientHTTPAPI[[]Server](d, req)
	if err != nil {
		return servers, errors.WithStack(err)
	}
	return servers, nil

}

// Healthcheck implements Client.
func (d *DefaultClient) Healthcheck(ctx context.Context) (Server, error) {
	url, err := url.JoinPath(d.baseURL, "api", "v1", "servers", "localhost")
	if err != nil {
		return Server{}, errors.WithStack(err)
	}
	req, err := d.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Server{}, errors.WithStack(err)
	}

	server, err := doDefaultClientHTTPAPI[Server](d, req)
	if err != nil {
		return Server{}, errors.WithStack(err)
	}
	return server, nil
}

func (d *DefaultClient) newRequest(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Add(powerDNSAPIKeyHeader, d.apiKey)
	return req, nil
}

func doDefaultClientHTTPAPI[T any](
	d *DefaultClient,
	req *http.Request,
) (T, error) {
	d.logger.DebugContext(req.Context(), "send http request", slog.String("url", req.URL.String()))

	var result T
	resp, err := d.client.Do(req)
	if err != nil {
		return result, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, errors.Newf("%s: must be 200 but got %d", req.URL.String(), resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, errors.WithStack(err)
	}
	return result, nil
}

func NewDefault(logger *slog.Logger, baseURL string, apiKey string, opts ...ClientOption) Client {
	client := http.Client{}

	for _, o := range opts {
		o(&client)
	}

	return &DefaultClient{logger: logger, baseURL: baseURL, client: client, apiKey: apiKey}
}

type ClientOption func(client *http.Client)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *http.Client) {
		client.Timeout = timeout
	}
}
