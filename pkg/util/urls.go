package util

import (
	"context"
	"net/http"
	"net/url"
)

func V1LoginDiscordRequest(ctx context.Context, base string, path string, successUrl string, state string) (*http.Request, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	u, err = u.Parse("/v1/login/discord")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("success_url", successUrl)
	q.Add("state", state)
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
