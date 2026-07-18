package router

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const shuttleQuery = `
query KakaoShuttle($after: LocalTime) {
  shuttle(input: {
    stops: [
      {name: "dormitory_o", limit: {order: 2, destination: 2}},
      {name: "shuttlecock_o", limit: {order: 2, destination: 2}},
      {name: "station", limit: {order: 2, destination: 2}},
      {name: "terminal", limit: {order: 2, destination: 2}},
      {name: "jungang_stn", limit: {order: 2, destination: 2}},
      {name: "shuttlecock_i", limit: {order: 2, destination: 2}}
    ],
    after: $after
  }) {
    stops {
      name
      timetable {
        destination {
          destination
          entries {
            time
            route { name tag }
          }
        }
      }
    }
  }
}`

type shuttleResult struct {
	Shuttle struct {
		Stops []shuttleStop `json:"stops"`
	} `json:"shuttle"`
}

type shuttleStop struct {
	Name      string `json:"name"`
	Timetable struct {
		Destination []shuttleDestination `json:"destination"`
	} `json:"timetable"`
}

type shuttleDestination struct {
	Destination string             `json:"destination"`
	Entries     []shuttleDeparture `json:"entries"`
}

type shuttleDeparture struct {
	Time  string       `json:"time"`
	Route shuttleRoute `json:"route"`
}

type shuttleRoute struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func GetShuttleMessage(ctx fiber.Ctx) error {
	if _, bindErr := bindSkillPayload(ctx); bindErr != nil {
		return badRequest(ctx, bindErr)
	}
	client, err := backendClient(ctx)
	if err != nil {
		return backendFailure(ctx, "shuttle", "셔틀", err)
	}

	currentTime := currentServiceTime()
	queryCtx, cancel := queryContext()
	defer cancel()
	var result shuttleResult
	if queryErr := client.Query(queryCtx, shuttleQuery, map[string]any{
		"after": currentTime.Format("15:04:05"),
	}, &result); queryErr != nil {
		return backendFailure(ctx, "shuttle", "셔틀", queryErr)
	}

	items := make([]any, 0, len(result.Shuttle.Stops))
	for _, stop := range result.Shuttle.Stops {
		items = append(items, shuttleCard(stop, currentTime.Format("15:04")))
	}
	if len(items) == 0 {
		return ctx.JSON(schema.TextResponse(
			"현재 예정된 셔틀 운행이 없어요.",
			navigationQuickReplies("셔틀"),
		))
	}

	return ctx.JSON(schema.CarouselResponse("itemCard", items, navigationQuickReplies("셔틀")))
}

func shuttleCard(stop shuttleStop, updatedAt string) schema.ItemCard {
	rows := make([]schema.ItemList, 0, len(stop.Timetable.Destination))
	for _, destination := range stop.Timetable.Destination {
		departures := make([]string, 0, len(destination.Entries))
		for _, entry := range destination.Entries {
			departures = append(departures, fmt.Sprintf(
				"%s %s",
				shortTime(entry.Time),
				shuttleRouteLabel(entry.Route),
			))
		}
		if len(departures) == 0 {
			departures = append(departures, "운행 없음")
		}
		rows = append(rows, schema.ItemList{
			Title:       shuttleDestinationName(destination.Destination),
			Description: strings.Join(departures, " · "),
		})
	}
	if len(rows) == 0 {
		rows = append(rows, schema.ItemList{Title: "운행 정보", Description: "현재 예정된 셔틀이 없어요"})
	}
	return schema.ItemCard{
		Title:       shuttleStopName(stop.Name),
		Description: updatedAt + " 기준 · 다음 운행",
		ItemList:    rows,
		Buttons:     []schema.Button{appButton("/shuttle?stop="+stop.Name, "전체 시간표")},
	}
}

func shuttleRouteLabel(route shuttleRoute) string {
	if route.Tag == "C" {
		return "순환"
	}
	if route.Name != "" {
		return route.Name
	}
	return "직행"
}

func shuttleStopName(stop string) string {
	switch stop {
	case "dormitory_o":
		return "기숙사"
	case "shuttlecock_o":
		return "셔틀콕"
	case "station":
		return "한대앞역"
	case "terminal":
		return "예술인아파트"
	case "jungang_stn":
		return "중앙역"
	case "shuttlecock_i":
		return "셔틀콕 건너편"
	default:
		return stop
	}
}

func shuttleDestinationName(destination string) string {
	switch destination {
	case "STATION":
		return "한대앞역 방면"
	case "TERMINAL":
		return "예술인 방면"
	case "JUNGANG":
		return "중앙역 방면"
	case "CAMPUS":
		return "캠퍼스 방면"
	default:
		return destination
	}
}
