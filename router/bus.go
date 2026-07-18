package router

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const busQuery = `
query KakaoBus {
  bus(input: [
    {route: 216000068, stop: 216000379, limit: 3},
    {route: 216000068, stop: 216000138, limit: 3},
    {route: 216000061, stop: 216000379, limit: 3},
    {route: 216000096, stop: 216000719, limit: 3},
    {route: 216000104, stop: 216000070, limit: 3},
    {route: 200000015, stop: 216000070, limit: 3},
    {route: 216000026, stop: 216000719, limit: 3},
    {route: 216000043, stop: 216000719, limit: 3},
    {route: 216000075, stop: 216000759, limit: 3}
  ]) {
    route { seq name }
    stop { seq name }
    arrival {
      stops
      seats
      minutes
      lowFloor
      isRealtime
      time
      arrivalTime
    }
  }
}`

const (
	route10Dash1 = 216000068
	route3102    = 216000061
	route3100N   = 216000096
	route7070    = 216000104
	route9090    = 200000015
	route3100    = 216000026
	route3101    = 216000043
	route50      = 216000075

	stopERICA             = 216000379
	stopSangnoksu         = 216000138
	stopMainGate          = 216000719
	stopHanyangUniversity = 216000070
	stopSeongpo           = 216000759

	maxBusArrivalPreview = 2
	busArrivalPartCount  = 4
)

type busResult struct {
	Bus []busRouteStop `json:"bus"`
}

type busRouteStop struct {
	Route   busRoute     `json:"route"`
	Stop    busStop      `json:"stop"`
	Arrival []busArrival `json:"arrival"`
}

type busRoute struct {
	Seq  int    `json:"seq"`
	Name string `json:"name"`
}

type busStop struct {
	Seq  int    `json:"seq"`
	Name string `json:"name"`
}

type busArrival struct {
	Stops       *int   `json:"stops"`
	Seats       *int   `json:"seats"`
	Minutes     *int   `json:"minutes"`
	LowFloor    *bool  `json:"lowFloor"`
	IsRealtime  bool   `json:"isRealtime"`
	Time        string `json:"time"`
	ArrivalTime string `json:"arrivalTime"`
}

type busKey struct {
	Route int
	Stop  int
}

type busCardSpec struct {
	Title  string
	Routes []busRouteSpec
}

type busRouteSpec struct {
	Key   busKey
	Label string
}

func GetBusMessage(ctx fiber.Ctx) error {
	if _, bindErr := bindSkillPayload(ctx); bindErr != nil {
		return badRequest(ctx, bindErr)
	}
	client, err := backendClient(ctx)
	if err != nil {
		return backendFailure(ctx, "bus", "버스", err)
	}

	queryCtx, cancel := queryContext()
	defer cancel()
	var result busResult
	if queryErr := client.Query(queryCtx, busQuery, nil, &result); queryErr != nil {
		return backendFailure(ctx, "bus", "버스", queryErr)
	}

	byKey := make(map[busKey]busRouteStop, len(result.Bus))
	for _, routeStop := range result.Bus {
		byKey[busKey{Route: routeStop.Route.Seq, Stop: routeStop.Stop.Seq}] = routeStop
	}
	cardSpecs := busCardSpecs()
	items := make([]any, 0, len(cardSpecs))
	for _, card := range cardSpecs {
		items = append(items, busListCard(card, byKey))
	}
	return ctx.JSON(schema.CarouselResponse("listCard", items, navigationQuickReplies("버스")))
}

func busCardSpecs() []busCardSpec {
	return []busCardSpec{
		{
			Title: "상록수역 방면",
			Routes: []busRouteSpec{
				{Key: busKey{Route: route10Dash1, Stop: stopERICA}, Label: "10-1 · ERICA"},
				{Key: busKey{Route: route10Dash1, Stop: stopSangnoksu}, Label: "10-1 · 상록수역"},
			},
		},
		{
			Title: "강남역 방면",
			Routes: []busRouteSpec{
				{Key: busKey{Route: route3102, Stop: stopERICA}, Label: "3102 · ERICA"},
				{Key: busKey{Route: route3100N, Stop: stopMainGate}, Label: "3100N · 정문"},
			},
		},
		{
			Title: "수원역 방면",
			Routes: []busRouteSpec{
				{Key: busKey{Route: route7070, Stop: stopHanyangUniversity}, Label: "7070 · 한양대입구"},
				{Key: busKey{Route: route9090, Stop: stopHanyangUniversity}, Label: "9090 · 한양대입구"},
			},
		},
		{
			Title: "군포·의왕 방면",
			Routes: []busRouteSpec{
				{Key: busKey{Route: route3100, Stop: stopMainGate}, Label: "3100 · 정문"},
				{Key: busKey{Route: route3101, Stop: stopMainGate}, Label: "3101 · 정문"},
			},
		},
		{
			Title: "광명역 방면",
			Routes: []busRouteSpec{
				{Key: busKey{Route: route50, Stop: stopSeongpo}, Label: "50 · 성포주공4단지"},
			},
		},
	}
}

func busListCard(spec busCardSpec, results map[busKey]busRouteStop) schema.ListCard {
	items := make([]schema.ListItem, 0, len(spec.Routes))
	for _, route := range spec.Routes {
		items = append(items, schema.ListItem{
			Title:       route.Label,
			Description: busArrivalDescription(results[route.Key].Arrival),
		})
	}
	return schema.ListCard{
		Header:  schema.ListItem{Title: spec.Title},
		Items:   items,
		Buttons: []schema.Button{appButton("/bus", "전체 버스")},
	}
}

func busArrivalDescription(arrivals []busArrival) string {
	if len(arrivals) == 0 {
		return noArrivalText
	}
	lines := make([]string, 0, min(len(arrivals), maxBusArrivalPreview))
	for index, arrival := range arrivals {
		if index == maxBusArrivalPreview {
			break
		}
		lines = append(lines, formatBusArrival(arrival))
	}
	return strings.Join(lines, " / ")
}

func formatBusArrival(arrival busArrival) string {
	if !arrival.IsRealtime {
		departure := arrival.ArrivalTime
		if departure == "" {
			departure = arrival.Time
		}
		if departure != "" {
			return shortTime(departure) + " 도착 예정"
		}
	}
	parts := make([]string, 0, busArrivalPartCount)
	if arrival.Minutes != nil {
		if *arrival.Minutes <= 0 {
			parts = append(parts, "곧 도착")
		} else {
			parts = append(parts, fmt.Sprintf("%d분", *arrival.Minutes))
		}
	}
	if arrival.Stops != nil {
		parts = append(parts, fmt.Sprintf("%d정거장", *arrival.Stops))
	}
	if arrival.Seats != nil && *arrival.Seats >= 0 {
		parts = append(parts, fmt.Sprintf("%d석", *arrival.Seats))
	}
	if arrival.LowFloor != nil && *arrival.LowFloor {
		parts = append(parts, "저상")
	}
	if len(parts) == 0 {
		parts = append(parts, "도착 정보 확인 중")
	}
	return strings.Join(parts, " · ")
}
