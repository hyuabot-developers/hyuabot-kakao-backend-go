package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/util"
	"sort"
	"strconv"
	"strings"
)

type BusSuwonArrival struct {
	RouteName string
	Arrival   model.StopRouteRealtimeArrival
}

func BusArrival(c *fiber.Ctx) error {
	cardList := make([]model.TextCard, 0)
	quickReplies := make([]model.QuickReply, 0)

	campus := util.GetBusArrival(216000379)
	mainGate := util.GetBusArrival(216000719)
	sangnoksuStation := util.GetBusArrival(216000138)
	seonganHighSchool := util.GetBusArrival(216000070)

	cityBusString := ""
	gangnamBusString := ""
	suwonBusString := ""
	gunpoBusString := ""
	nightBusString := ""
	for _, route := range campus.Route {
		if route.Name == "10-1" {
			cityBusString += "10-1번 (학교 → 상록수역)\n"
			for _, arrival := range route.Realtime {
				cityBusString += strconv.Itoa(arrival.Time) + "분 후 도착\n"
			}
			for i, timetable := range route.Start.Timetable {
				cityBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		} else if route.Name == "3102" {
			gangnamBusString += "3102번\n"
			for _, arrival := range route.Realtime {
				gangnamBusString += strconv.Itoa(arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Seat) + "석)\n"
			}
			for i, timetable := range route.Start.Timetable {
				gangnamBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		}
	}

	for _, route := range mainGate.Route {
		if route.Name == "3100" {
			gunpoBusString += "\n3100번\n"
			for _, arrival := range route.Realtime {
				gunpoBusString += strconv.Itoa(arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Seat) + "석)\n"
			}
			for i, timetable := range route.Start.Timetable {
				gunpoBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		} else if route.Name == "3101" {
			gunpoBusString += "\n3101번\n"
			for _, arrival := range route.Realtime {
				gunpoBusString += strconv.Itoa(arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Seat) + "석)\n"
			}
			for i, timetable := range route.Start.Timetable {
				gunpoBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		} else if route.Name == "3100N" {
			nightBusString += "\n3100N번\n"
			for _, arrival := range route.Realtime {
				nightBusString += strconv.Itoa(arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Seat) + "석)\n"
			}
			for i, timetable := range route.Start.Timetable {
				nightBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		} else if route.Name == "707-1" {
			suwonBusString += "\n707-1번\n"
			for _, arrival := range route.Realtime {
				suwonBusString += strconv.Itoa(arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Seat) + "석)\n"
			}
			for i, timetable := range route.Start.Timetable {
				suwonBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2-len(route.Realtime) {
					break
				}
			}
		}
	}

	for _, route := range sangnoksuStation.Route {
		if route.Name == "10-1" {
			cityBusString += "\n10-1번 (상록수역 → 학교)\n"
			for i, timetable := range route.Start.Timetable {
				cityBusString += stripTimetable(timetable) + "분 출발\n"
				if i >= 2 {
					break
				}
			}
		}
	}

	suwonRealtimeArrival := make([]BusSuwonArrival, 0)
	for _, route := range seonganHighSchool.Route {
		if route.Name == "110" || route.Name == "707" || route.Name == "909" {
			for _, arrival := range route.Realtime {
				suwonRealtimeArrival = append(suwonRealtimeArrival, BusSuwonArrival{RouteName: route.Name, Arrival: arrival})
			}
		}
	}
	sort.Slice(suwonRealtimeArrival, func(i, j int) bool {
		return suwonRealtimeArrival[i].Arrival.Time < suwonRealtimeArrival[j].Arrival.Time
	})
	suwonBusString += "\n110/707/909번\n"
	for i, arrival := range suwonRealtimeArrival {
		suwonBusString += "(" + arrival.RouteName + "번) " + strconv.Itoa(arrival.Arrival.Time) + "분 후 도착(" + strconv.Itoa(arrival.Arrival.Seat) + "석)\n"
		if i >= 4 {
			break
		}
	}

	cityBusCard := model.TextCard{
		Title:       "시내버스",
		Description: cityBusString,
		Buttons:     make([]model.CardButton, 0),
	}
	gangnamBusCard := model.TextCard{
		Title:       "강남역 방면",
		Description: gangnamBusString,
		Buttons:     make([]model.CardButton, 0),
	}
	gunpoBusCard := model.TextCard{
		Title:       "군포, 의왕 방면",
		Description: gunpoBusString,
		Buttons:     make([]model.CardButton, 0),
	}
	nightBusCard := model.TextCard{
		Title:       "강남역 방면(심야)",
		Description: nightBusString,
		Buttons:     make([]model.CardButton, 0),
	}
	suwonBusCard := model.TextCard{
		Title:       "수원역 방면",
		Description: suwonBusString,
		Buttons:     make([]model.CardButton, 0),
	}
	cardList = append(cardList, cityBusCard)
	cardList = append(cardList, gangnamBusCard)
	cardList = append(cardList, nightBusCard)
	cardList = append(cardList, suwonBusCard)
	cardList = append(cardList, gunpoBusCard)
	response := util.SetResponse(
		util.SetTemplate([]model.Components{util.SetBasicCardCarousel(cardList)}, quickReplies))
	return c.JSON(response)
}

func stripTimetable(timetable string) string {
	return strings.Replace(timetable[:len(timetable)-3], ":", "시 ", 1)
}
