package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func JsonFetch[T any](ctx context.Context, url url.URL, method string) (*T, error) {
	if len(url.String()) == 0 {
		return nil, fmt.Errorf("invalid url provided")
	}

	if len(method) == 0 {
		return nil, fmt.Errorf("invalid method provided")
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to form request:: %w", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch:: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusPermanentRedirect {
		return nil, fmt.Errorf("server returned a non-successful response code ::%s", res.Status)
	}

	rawBytes, err := io.ReadAll(res.Body)

	response := new(T)

	if err := json.Unmarshal(rawBytes, response); err != nil {
		return nil, fmt.Errorf("failed to parse response:: %w", err)
	}

	return response, nil
}
