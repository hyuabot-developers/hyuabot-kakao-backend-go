package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const subwayQuery = `
query KakaoSubway($weekday: String!) {
  subway(input: {keys: [
    {stationID: "K449", direction: ["up", "down"], weekdays: [$weekday], limit: 4},
    {stationID: "K251", direction: ["up", "down"], weekdays: [$weekday], limit: 4}
  ]}) {
    stationID
    name
    route { name }
    arrival {
      direction
      entries {
        minutes
        isRealtime
        location
        stops
        isExpress
        isLast
        terminal { name }
      }
    }
  }
}`

const maxSubwayArrivalPreview = 2

type subwayResult struct {
	Subway []subwayStation `json:"subway"`
}

type subwayStation struct {
	StationID string               `json:"stationID"`
	Name      string               `json:"name"`
	Route     subwayRoute          `json:"route"`
	Arrival   []subwayArrivalGroup `json:"arrival"`
}

type subwayRoute struct {
	Name string `json:"name"`
}

type subwayArrivalGroup struct {
	Direction string          `json:"direction"`
	Entries   []subwayArrival `json:"entries"`
}

type subwayArrival struct {
	Minutes    int            `json:"minutes"`
	IsRealtime bool           `json:"isRealtime"`
	Location   string         `json:"location"`
	Stops      *int           `json:"stops"`
	IsExpress  *bool          `json:"isExpress"`
	IsLast     *bool          `json:"isLast"`
	Terminal   subwayTerminal `json:"terminal"`
}

type subwayTerminal struct {
	Name string `json:"name"`
}

func GetSubwayMessage(ctx fiber.Ctx) error {
	if _, bindErr := bindSkillPayload(ctx); bindErr != nil {
		return badRequest(ctx, bindErr)
	}
	client, err := backendClient(ctx)
	if err != nil {
		return backendFailure(ctx, "subway", "지하철", err)
	}

	currentTime := currentServiceTime()
	weekday := "weekdays"
	if currentTime.Weekday() == time.Saturday || currentTime.Weekday() == time.Sunday {
		weekday = "weekends"
	}
	queryCtx, cancel := queryContext()
	defer cancel()
	var result subwayResult
	if queryErr := client.Query(queryCtx, subwayQuery, map[string]any{"weekday": weekday}, &result); queryErr != nil {
		return backendFailure(ctx, "subway", "지하철", queryErr)
	}

	items := make([]any, 0, len(result.Subway))
	for _, station := range result.Subway {
		items = append(items, subwayCard(station, currentTime.Format("15:04")))
	}
	if len(items) == 0 {
		return ctx.JSON(schema.TextResponse(
			"현재 지하철 도착 정보가 없어요.",
			navigationQuickReplies("지하철"),
		))
	}
	return ctx.JSON(schema.CarouselResponse("itemCard", items, navigationQuickReplies("지하철")))
}

func subwayCard(station subwayStation, updatedAt string) schema.ItemCard {
	rows := make([]schema.ItemList, 0, len(station.Arrival))
	for _, group := range station.Arrival {
		rows = append(rows, schema.ItemList{
			Title:       subwayDirectionName(station.StationID, group.Direction),
			Description: subwayArrivalDescription(group.Entries),
		})
	}
	if len(rows) == 0 {
		rows = append(rows, schema.ItemList{Title: "도착 정보", Description: noArrivalText})
	}
	title := station.Route.Name
	if title == "" {
		title = station.Name
	}
	return schema.ItemCard{
		Title:       title + " · " + station.Name,
		Description: updatedAt + " 기준",
		ItemList:    rows,
		Buttons:     []schema.Button{appButton("/subway", "전체 지하철")},
	}
}

func subwayDirectionName(stationID string, direction string) string {
	switch stationID + ":" + direction {
	case "K449:up":
		return "당고개 방면"
	case "K449:down":
		return "오이도 방면"
	case "K251:up":
		return "청량리 방면"
	case "K251:down":
		return "인천 방면"
	default:
		return direction + " 방면"
	}
}

func subwayArrivalDescription(arrivals []subwayArrival) string {
	if len(arrivals) == 0 {
		return noArrivalText
	}
	lines := make([]string, 0, min(len(arrivals), maxSubwayArrivalPreview))
	for index, arrival := range arrivals {
		if index == maxSubwayArrivalPreview {
			break
		}
		parts := []string{fmt.Sprintf("%d분", arrival.Minutes)}
		if arrival.Terminal.Name != "" {
			parts = append(parts, arrival.Terminal.Name+"행")
		}
		if arrival.Location != "" {
			parts = append(parts, arrival.Location)
		} else if arrival.Stops != nil {
			parts = append(parts, fmt.Sprintf("%d정거장 전", *arrival.Stops))
		}
		if arrival.IsExpress != nil && *arrival.IsExpress {
			parts = append(parts, "급행")
		}
		if !arrival.IsRealtime {
			parts = append(parts, "시간표")
		}
		lines = append(lines, strings.Join(parts, " · "))
	}
	return strings.Join(lines, " / ")
}
