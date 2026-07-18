package backend_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/backend"
)

func TestHTTPClientQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if request.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %s", request.Header.Get("Content-Type"))
		}
		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if body["query"] != "query Test { value }" {
			t.Errorf("query = %v", body["query"])
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(writer, `{"data":{"value":"ok"}}`)
	}))
	t.Cleanup(server.Close)

	client := backend.NewClient(server.URL, server.Client())
	var result struct {
		Value string `json:"value"`
	}
	if err := client.Query(context.Background(), "query Test { value }", nil, &result); err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if result.Value != "ok" {
		t.Errorf("value = %q, want ok", result.Value)
	}
}

func TestHTTPClientQueryErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		body       string
		want       string
	}{
		{name: "HTTP status", statusCode: http.StatusBadGateway, body: `{}`, want: "HTTP 502"},
		{
			name:       "GraphQL error",
			statusCode: http.StatusOK,
			body:       `{"errors":[{"message":"invalid input"}]}`,
			want:       "invalid input",
		},
		{name: "missing data", statusCode: http.StatusOK, body: `{"data":null}`, want: "did not contain data"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(test.statusCode)
				_, _ = io.WriteString(writer, test.body)
			}))
			t.Cleanup(server.Close)

			client := backend.NewClient(server.URL, server.Client())
			var result map[string]any
			err := client.Query(context.Background(), "query Test { value }", nil, &result)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Query() error = %v, want containing %q", err, test.want)
			}
		})
	}
}
