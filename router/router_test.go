package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/router"
)

type fakeClient struct {
	data      any
	err       error
	query     string
	variables map[string]any
}

func (client *fakeClient) Query(
	_ context.Context,
	query string,
	variables map[string]any,
	result any,
) error {
	client.query = query
	client.variables = variables
	if client.err != nil {
		return client.err
	}
	data, err := json.Marshal(client.data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func TestSkillRoutesUseCurrentGraphQLSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		path         string
		handler      fiber.Handler
		data         any
		queryField   string
		carouselType string
	}{
		{
			name:         "shuttle",
			path:         "/shuttle",
			handler:      router.GetShuttleMessage,
			queryField:   "shuttle(input:",
			carouselType: "itemCard",
			data: map[string]any{"shuttle": map[string]any{"stops": []any{
				map[string]any{
					"name": "dormitory_o",
					"timetable": map[string]any{"destination": []any{
						map[string]any{
							"destination": "STATION",
							"entries": []any{
								map[string]any{
									"time":  "10:30:00",
									"route": map[string]any{"name": "직행", "tag": "D"},
								},
							},
						},
					}},
				},
			}}},
		},
		{
			name:         "bus",
			path:         "/bus",
			handler:      router.GetBusMessage,
			queryField:   "bus(input:",
			carouselType: "listCard",
			data: map[string]any{"bus": []any{
				map[string]any{
					"route": map[string]any{"seq": 216000068, "name": "10-1"},
					"stop":  map[string]any{"seq": 216000379, "name": "ERICA컨벤션센터"},
					"arrival": []any{
						map[string]any{
							"stops": 2, "seats": 17, "minutes": 3,
							"lowFloor": true, "isRealtime": true,
						},
					},
				},
			}},
		},
		{
			name:         "cafeteria",
			path:         "/cafeteria",
			handler:      router.GetCafeteriaMessage,
			queryField:   "cafeteria(input:",
			carouselType: "listCard",
			data: map[string]any{"cafeteria": []any{
				map[string]any{
					"seq": 12, "name": "학생식당",
					"menus": []any{
						map[string]any{"type": "중식", "food": "제육덮밥", "price": "5,500원"},
					},
				},
			}},
		},
		{
			name:         "subway",
			path:         "/subway",
			handler:      router.GetSubwayMessage,
			queryField:   "subway(input:",
			carouselType: "itemCard",
			data: map[string]any{"subway": []any{
				map[string]any{
					"stationID": "K449", "name": "한대앞", "route": map[string]any{"name": "4호선"},
					"arrival": []any{
						map[string]any{
							"direction": "up",
							"entries": []any{
								map[string]any{
									"minutes": 3, "isRealtime": true, "location": "전역 출발",
									"terminal": map[string]any{"name": "당고개"},
								},
							},
						},
					},
				},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &fakeClient{data: test.data}
			app := newTestApp(client)
			app.Post(test.path, test.handler)

			response := performRequest(t, app, http.MethodPost, test.path, skillPayload("중식 메뉴"))
			if response.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", response.StatusCode)
			}
			body := decodeResponse(t, response)
			assertCarouselResponse(t, body, test.carouselType)
			if !strings.Contains(client.query, test.queryField) {
				t.Errorf("query does not contain %q: %s", test.queryField, client.query)
			}
			for _, legacyField := range []string{"startStr", "dateStr", "id_:", "groupedTimetable"} {
				if strings.Contains(client.query, legacyField) {
					t.Errorf("query contains legacy field %q", legacyField)
				}
			}
		})
	}
}

func TestSkillRouteReturnsFriendlyFallback(t *testing.T) {
	t.Parallel()

	client := &fakeClient{err: errors.New("backend unavailable")}
	app := newTestApp(client)
	app.Post("/shuttle", router.GetShuttleMessage)

	response := performRequest(t, app, http.MethodPost, "/shuttle", skillPayload("셔틀"))
	body := decodeResponse(t, response)
	outputs := body["template"].(map[string]any)["outputs"].([]any)
	text := outputs[0].(map[string]any)["simpleText"].(map[string]any)["text"].(string)
	if !strings.Contains(text, "다시 시도") {
		t.Errorf("fallback text = %q", text)
	}
}

func TestHealthcheckSupportsReadinessProbe(t *testing.T) {
	t.Parallel()

	client := &fakeClient{data: map[string]any{"__typename": "Query"}}
	app := newTestApp(client)
	app.Get("/healthcheck", router.GetHealthCheckMessage)

	response := performRequest(t, app, http.MethodGet, "/healthcheck", nil)
	body := decodeResponse(t, response)
	if body["status"] != "ok" {
		t.Errorf("status body = %v", body)
	}
	if !strings.Contains(client.query, "__typename") {
		t.Errorf("healthcheck query = %q", client.query)
	}
}

func newTestApp(client *fakeClient) *fiber.App {
	app := fiber.New()
	app.Use(func(ctx fiber.Ctx) error {
		ctx.Locals(router.BackendClientLocal, client)
		return ctx.Next()
	})
	return app
}

func skillPayload(utterance string) []byte {
	payload := map[string]any{
		"action":      map[string]any{"params": map[string]string{}},
		"userRequest": map[string]any{"utterance": utterance},
	}
	body, _ := json.Marshal(payload)
	return body
}

func performRequest(t *testing.T, app *fiber.App, method string, path string, body []byte) *http.Response {
	t.Helper()
	request, err := http.NewRequestWithContext(context.Background(), method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("perform request: %v", err)
	}
	return response
}

func decodeResponse(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Errorf("close response: %v", err)
		}
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	var decoded map[string]any
	if decodeErr := json.Unmarshal(body, &decoded); decodeErr != nil {
		t.Fatalf("decode response %q: %v", body, decodeErr)
	}
	return decoded
}

func assertCarouselResponse(t *testing.T, body map[string]any, carouselType string) {
	t.Helper()
	if body["version"] != "2.0" {
		t.Errorf("version = %v", body["version"])
	}
	template := body["template"].(map[string]any)
	outputs := template["outputs"].([]any)
	carousel := outputs[0].(map[string]any)["carousel"].(map[string]any)
	if carousel["type"] != carouselType {
		t.Errorf("carousel type = %v, want %s", carousel["type"], carouselType)
	}
	quickReplies := template["quickReplies"].([]any)
	if len(quickReplies) > 10 {
		t.Errorf("quick replies = %d, want <= 10", len(quickReplies))
	}
}
