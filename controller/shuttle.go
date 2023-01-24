package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/util"
	"sort"
	"strconv"
	"strings"
)

type ShuttleArrivalItem struct {
	Tag         string
	ArrivalTime int
}

func ShuttleArrival(c *fiber.Ctx) error {
	arrivalResponse := util.GetShuttleArrival()

	dormitoryStation := make([]ShuttleArrivalItem, 0)
	dormitoryTerminal := make([]ShuttleArrivalItem, 0)
	dormitoryJungang := make([]ShuttleArrivalItem, 0)

	shuttlecockOutStation := make([]ShuttleArrivalItem, 0)
	shuttlecockOutTerminal := make([]ShuttleArrivalItem, 0)
	shuttlecockOutJungang := make([]ShuttleArrivalItem, 0)

	stationCampus := make([]ShuttleArrivalItem, 0)
	stationTerminal := make([]ShuttleArrivalItem, 0)
	stationJungang := make([]ShuttleArrivalItem, 0)

	terminalCampus := make([]ShuttleArrivalItem, 0)

	jungangCampus := make([]ShuttleArrivalItem, 0)

	shuttlecockIn := make([]ShuttleArrivalItem, 0)

	for _, stop := range arrivalResponse.Stop {
		if stop.Name == "dormitory_o" {
			for _, route := range stop.Route {
				if route.Tag == "DH" {
					for _, arrival := range route.Arrival {
						dormitoryStation = append(dormitoryStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "DY" {
					for _, arrival := range route.Arrival {
						dormitoryTerminal = append(dormitoryTerminal, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "DJ" {
					for _, arrival := range route.Arrival {
						dormitoryStation = append(dormitoryStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						dormitoryJungang = append(dormitoryJungang, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "C" {
					for _, arrival := range route.Arrival {
						dormitoryStation = append(dormitoryStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						dormitoryTerminal = append(dormitoryTerminal, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				}
			}
		} else if stop.Name == "shuttlecock_o" {
			for _, route := range stop.Route {
				if route.Tag == "DH" {
					for _, arrival := range route.Arrival {
						shuttlecockOutStation = append(shuttlecockOutStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "DY" {
					for _, arrival := range route.Arrival {
						shuttlecockOutTerminal = append(shuttlecockOutTerminal, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "DJ" {
					for _, arrival := range route.Arrival {
						shuttlecockOutStation = append(shuttlecockOutStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						shuttlecockOutJungang = append(shuttlecockOutJungang, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "C" {
					for _, arrival := range route.Arrival {
						shuttlecockOutStation = append(shuttlecockOutStation, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						shuttlecockOutTerminal = append(shuttlecockOutTerminal, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				}
			}
		} else if stop.Name == "station" {
			for _, route := range stop.Route {
				if route.Tag == "C" {
					for _, arrival := range route.Arrival {
						stationCampus = append(stationCampus, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						stationTerminal = append(stationTerminal, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else if route.Tag == "DJ" {
					for _, arrival := range route.Arrival {
						stationCampus = append(stationCampus, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
						stationJungang = append(stationJungang, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				} else {
					for _, arrival := range route.Arrival {
						stationCampus = append(stationCampus, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				}
			}
		} else if stop.Name == "terminal" {
			for _, route := range stop.Route {
				for _, arrival := range route.Arrival {
					terminalCampus = append(terminalCampus, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
				}
			}
		} else if stop.Name == "jungang_stn" {
			for _, route := range stop.Route {
				for _, arrival := range route.Arrival {
					jungangCampus = append(jungangCampus, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
				}
			}
		} else if stop.Name == "shuttlecock_i" {
			for _, route := range stop.Route {
				if strings.HasSuffix(route.Name, "D") {
					for _, arrival := range route.Arrival {
						shuttlecockIn = append(shuttlecockIn, ShuttleArrivalItem{Tag: route.Tag, ArrivalTime: arrival})
					}
				}
			}
		}
	}

	dormitoryString := ""
	dormitoryString += "한대앞 방면\n"

	if len(dormitoryStation) == 0 {
		dormitoryString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(dormitoryStation) {
			dormitoryString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	dormitoryString += "예술인 방면\n"
	if len(dormitoryTerminal) == 0 {
		dormitoryString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(dormitoryTerminal) {
			dormitoryString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	if len(dormitoryJungang) > 0 {
		dormitoryString += "중앙역 방면\n"
		for i, item := range sortShuttleArrivalItem(dormitoryJungang) {
			dormitoryString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	dormitoryCard := model.TextCard{
		Title:       "기숙사",
		Description: strings.TrimRight(dormitoryString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	shuttlecockOutString := ""
	shuttlecockOutString += "한대앞 방면\n"
	if len(shuttlecockOutStation) == 0 {
		shuttlecockOutString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(shuttlecockOutStation) {
			shuttlecockOutString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	shuttlecockOutString += "예술인 방면\n"
	if len(shuttlecockOutTerminal) == 0 {
		shuttlecockOutString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(shuttlecockOutTerminal) {
			shuttlecockOutString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	if len(shuttlecockOutJungang) > 0 {
		shuttlecockOutString += "중앙역 방면\n"
		for i, item := range sortShuttleArrivalItem(shuttlecockOutTerminal) {
			shuttlecockOutString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	shuttlecockOutCard := model.TextCard{
		Title:       "셔틀콕",
		Description: strings.Trim(shuttlecockOutString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	stationString := ""
	stationString += "캠퍼스 방면\n"
	if len(stationCampus) == 0 {
		stationString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(stationCampus) {
			stationString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	if len(stationJungang) > 0 {
		stationString += "중앙역 방면\n"
		for i, item := range sortShuttleArrivalItem(stationJungang) {
			stationString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	stationString += "예술인 방면\n"
	if len(stationTerminal) == 0 {
		stationString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(stationTerminal) {
			stationString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	stationCard := model.TextCard{
		Title:       "한대앞역",
		Description: strings.Trim(stationString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	terminalString := ""
	terminalString += "캠퍼스 방면\n"
	if len(terminalCampus) == 0 {
		terminalString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(terminalCampus) {
			terminalString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 6 {
				break
			}
		}
	}
	terminalCard := model.TextCard{
		Title:       "예술인",
		Description: strings.Trim(terminalString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	jungangString := ""
	jungangString += "캠퍼스 방면\n"
	if len(jungangCampus) == 0 {
		jungangString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(jungangCampus) {
			jungangString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 6 {
				break
			}
		}
	}
	jungangCard := model.TextCard{
		Title:       "중앙역",
		Description: strings.Trim(jungangString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	shuttlecockInString := ""
	shuttlecockInString += "기숙사 방면\n"
	if len(shuttlecockIn) == 0 {
		shuttlecockInString += "운행 종료\n"
	} else {
		for i, item := range sortShuttleArrivalItem(shuttlecockIn) {
			shuttlecockInString += strconv.Itoa(item.ArrivalTime) + "분 후 도착(" + getHeadingString(item.Tag) + ")\n"
			if i == 2 {
				break
			}
		}
	}
	shuttlecockInCard := model.TextCard{
		Title:       "셔틀콕 건너편",
		Description: strings.Trim(shuttlecockInString, "\n"),
		Buttons:     make([]model.CardButton, 0),
	}

	cardList := make([]model.TextCard, 0)
	quickReplies := make([]model.QuickReply, 0)
	cardList = append(cardList, dormitoryCard)
	cardList = append(cardList, shuttlecockOutCard)
	cardList = append(cardList, stationCard)
	cardList = append(cardList, terminalCard)
	cardList = append(cardList, jungangCard)
	cardList = append(cardList, shuttlecockInCard)

	quickReplies = append(quickReplies, model.QuickReply{Action: "block", Label: "🏘️ 기숙사", MessageText: "🏘️ 기숙사", BlockID: "5ebf702e7a9c4b000105fb25"})
	quickReplies = append(quickReplies, model.QuickReply{Action: "block", Label: "🏫  셔틀콕", MessageText: "🏫 셔틀콕", BlockID: "5ebf702e7a9c4b000105fb25"})
	quickReplies = append(quickReplies, model.QuickReply{Action: "block", Label: "🚆 한대앞역", MessageText: "🚆 한대앞역", BlockID: "5ebf702e7a9c4b000105fb25"})
	quickReplies = append(quickReplies, model.QuickReply{Action: "block", Label: "🚍 예술인A", MessageText: "🚍 예술인A", BlockID: "5ebf702e7a9c4b000105fb25"})
	quickReplies = append(quickReplies, model.QuickReply{Action: "block", Label: "🏫 셔틀콕 건너편", MessageText: "🏫 셔틀콕 건너편", BlockID: "5ebf702e7a9c4b000105fb25"})
	response := util.SetResponse(
		util.SetTemplate([]model.Components{util.SetBasicCardCarousel(cardList)}, quickReplies))
	return c.JSON(response)
}

func sortShuttleArrivalItem(arrivalItems []ShuttleArrivalItem) []ShuttleArrivalItem {
	sort.Slice(arrivalItems, func(i, j int) bool {
		return arrivalItems[i].ArrivalTime < arrivalItems[j].ArrivalTime
	})
	return arrivalItems
}

func getHeadingString(routeTag string) string {
	if routeTag == "C" {
		return "순환"
	}
	return "직행"
}
