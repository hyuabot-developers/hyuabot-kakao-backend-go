package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const maxResponseSize = 2 << 20

type Client interface {
	Query(ctx context.Context, query string, variables map[string]any, result any) error
}

type HTTPClient struct {
	endpoint   string
	httpClient *http.Client
}

type request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type response struct {
	Data   json.RawMessage `json:"data"`
	Errors []queryError    `json:"errors"`
}

type queryError struct {
	Message string `json:"message"`
}

func NewClient(endpoint string, httpClient *http.Client) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &HTTPClient{endpoint: endpoint, httpClient: httpClient}
}

func (client *HTTPClient) Query(
	ctx context.Context,
	query string,
	variables map[string]any,
	result any,
) error {
	body, err := json.Marshal(request{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("encode GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create GraphQL request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send GraphQL request: %w", err)
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(res.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("read GraphQL response: %w", err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("GraphQL server returned HTTP %d", res.StatusCode)
	}

	var envelope response
	if decodeErr := json.Unmarshal(responseBody, &envelope); decodeErr != nil {
		return fmt.Errorf("decode GraphQL response: %w", decodeErr)
	}
	if len(envelope.Errors) > 0 {
		messages := make([]error, 0, len(envelope.Errors))
		for _, queryError := range envelope.Errors {
			messages = append(messages, errors.New(queryError.Message))
		}
		return fmt.Errorf("GraphQL query failed: %w", errors.Join(messages...))
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return errors.New("GraphQL response did not contain data")
	}
	if decodeErr := json.Unmarshal(envelope.Data, result); decodeErr != nil {
		return fmt.Errorf("decode GraphQL data: %w", decodeErr)
	}
	return nil
}
