package router

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hasura/go-graphql-client"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

type BusStop struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
	Routes    []BusRoute
}

type BusRoute struct {
	Info      BusRouteInfo
	Timetable []BusTimetable
	Realtime  []BusRealtime
}

type BusRouteInfo struct {
	ID   int
	Name string
}

type BusTimetable struct {
	Weekdays string
	Time     string
}

type BusRealtime struct {
	Sequence int
	Stop     int
	Time     float64
	Seat     int
	LowFloor bool
}

type MergedBusRealtime struct {
	Name string
	Stop int
	Time int
	Seat int
}

const arrivalSectionLength = 3

func QueryBusDepartureData(ctx fiber.Ctx) []BusStop {
	// GraphQL Client and check API server status
	client, loaded := ctx.Locals("graphQLClient").(*graphql.Client)
	if !loaded {
		panic("GraphQL client not found")
	}
	// Get current datetime
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	currentTime := time.Now().In(location)
	// Query Shuttle Timetable
	var query struct {
		Bus []BusStop `graphql:"bus(id_: [216000138, 216000759, 216000381, 216000117, 216000379, 216000383, 216000070, 216000719, 213000487], startStr: $time)"`
	}
	variables := map[string]interface{}{
		"time": currentTime.Format("15:04:03"),
	}
	queryError := client.Query(context.Background(), &query, variables)
	if queryError != nil {
		panic(queryError)
	}
	return query.Bus
}

func GenerateBusSectionText(header string, result BusRoute) string {
	cardText := ""
	cardText += header
	for _, realtime := range result.Realtime {
		if realtime.Seat >= 0 {
			cardText += fmt.Sprintf("%d분 후 도착(%d석)\n", int(realtime.Time), realtime.Seat)
		} else {
			cardText += fmt.Sprintf("%d분 후 도착\n", int(realtime.Time))
		}
	}
	if len(result.Realtime) < arrivalSectionLength {
		for index, timetable := range result.Timetable {
			if index < arrivalSectionLength-len(result.Realtime) {
				cardText += fmt.Sprintf("%s 출발\n", strings.TrimSuffix(timetable.Time, ":00"))
			}
		}
	}
	if len(result.Realtime) == 0 && len(result.Timetable) == 0 {
		cardText += noArrivalText
	}
	return cardText
}

func GenerateMergedBusSectionText(header string, result []MergedBusRealtime) string {
	cardText := ""
	cardText += header
	for index, timetable := range result {
		cardText += fmt.Sprintf("%s %d분 후 도착(%d석)\n", timetable.Name, timetable.Time, timetable.Seat)
		if index == arrivalSectionLength-1 {
			break
		}
	}
	return cardText
}

func GetSangnoksuStationText(result map[int]map[int]BusRoute) string {
	campus := result[216000379][216000068]
	sangnoksu := result[216000138][216000069]
	cardText := ""
	cardText += GenerateBusSectionText("10-1 (ERICA)\n", campus)
	cardText += GenerateBusSectionText("\n10-1 (상록수역)\n", sangnoksu)
	return cardText
}

func GetGangnamStationText(result map[int]map[int]BusRoute) string {
	bus3102 := result[216000381][216000068]
	bus3100N := result[216000719][216000096]
	cardText := ""
	cardText += GenerateBusSectionText("3102 (ERICA)\n", bus3102)
	cardText += GenerateBusSectionText("\n3100N (한양대정문)\n", bus3100N)
	return cardText
}

func GetSuwonStationText(result map[int]map[int]BusRoute) string {
	bus7071 := result[216000719][216000070]
	bus7070 := result[216000070][216000104]
	bus110 := result[216000070][217000014]
	bus9090 := result[216000070][200000015]
	cardText := ""
	cardText += GenerateBusSectionText("707-1 (한양대정문)\n", bus7071)
	mergedBusRealtime := make([]MergedBusRealtime, 0)
	for _, realtime := range bus7070.Realtime {
		mergedBusRealtime = append(mergedBusRealtime, MergedBusRealtime{
			Name: bus7070.Info.Name,
			Stop: realtime.Stop,
			Time: int(realtime.Time),
			Seat: realtime.Seat,
		})
	}
	for _, realtime := range bus110.Realtime {
		mergedBusRealtime = append(mergedBusRealtime, MergedBusRealtime{
			Name: bus110.Info.Name,
			Stop: realtime.Stop,
			Time: int(realtime.Time),
			Seat: realtime.Seat,
		})
	}
	for _, realtime := range bus9090.Realtime {
		mergedBusRealtime = append(mergedBusRealtime, MergedBusRealtime{
			Name: bus9090.Info.Name,
			Stop: realtime.Stop,
			Time: int(realtime.Time),
			Seat: realtime.Seat,
		})
	}
	// Sort by time
	sort.Slice(mergedBusRealtime, func(i, j int) bool {
		return mergedBusRealtime[i].Time < mergedBusRealtime[j].Time
	})
	cardText += GenerateMergedBusSectionText("\n기타 (성안고)\n", mergedBusRealtime)
	return cardText
}

func GetGunpoText(result map[int]map[int]BusRoute) string {
	bus3100 := result[216000719][216000026]
	bus3101 := result[216000719][216000043]
	cardText := ""
	cardText += GenerateBusSectionText("3100 (한양대정문)\n", bus3100)
	cardText += GenerateBusSectionText("\n3101 (한양대정문)\n", bus3101)
	return cardText
}

func GetGwangmyeongText(result map[int]map[int]BusRoute) string {
	bus50 := result[216000759][216000075]
	cardText := ""
	cardText += GenerateBusSectionText("50 (성포주공4단지)\n", bus50)
	return cardText
}

func GetBusMessage(ctx fiber.Ctx) error {
	body := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Group shuttle timetable by stop and destination
	result := QueryBusDepartureData(ctx)
	resultMap := make(map[int]map[int]BusRoute)
	for _, busStop := range result {
		if resultMap[busStop.ID] == nil {
			resultMap[busStop.ID] = make(map[int]BusRoute)
		}
		for _, busRoute := range busStop.Routes {
			resultMap[busStop.ID][busRoute.Info.ID] = busRoute
		}
	}
	sangnoksuText := GetSangnoksuStationText(resultMap)
	gangnamText := GetGangnamStationText(resultMap)
	suwonText := GetSuwonStationText(resultMap)
	gunpoText := GetGunpoText(resultMap)
	gwangMyeongText := GetGwangmyeongText(resultMap)
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.Carousel{
					Content: schema.CarouselContent{
						Type: "textCard",
						Items: []schema.Component{
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       "상록수역",
									Description: strings.Trim(sangnoksuText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       "강남역",
									Description: strings.Trim(gangnamText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       "수원역",
									Description: strings.Trim(suwonText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       "군포/의왕",
									Description: strings.Trim(gunpoText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       "광명역",
									Description: strings.Trim(gwangMyeongText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
						},
					},
				},
			},
			QuickReplies: []schema.QuickReply{
				{
					Label:       "휴아봇 앱 설치",
					Action:      "block",
					MessageText: "휴아봇 앱 설치",
					BlockID:     "6077ca2de2039a2ba38c755f",
					Extra:       map[string]string{},
				},
			},
		},
	}
	return ctx.JSON(response)
}
