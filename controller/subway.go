package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/util"
	"strconv"
	"strings"
)

func Subway(c *fiber.Ctx) error {
	line4 := util.GetSubwayArrival("K449")
	lineSuinBundang := util.GetSubwayArrival("K251")

	var cardList []model.TextCard
	message := ""

	message += "서울 방면\n"
	for i, item := range line4.Realtime.Up {
		message += item.Terminal + "행 " + strconv.Itoa(item.Time) + "분 후 도착 (" + item.Location + ")\n"
		if i >= 3 {
			break
		}
	}
	message += "\n오이도 방면\n"
	for i, item := range line4.Realtime.Down {
		message += item.Terminal + "행 " + strconv.Itoa(item.Time) + "분 후 도착 (" + item.Location + ")\n"
		if i >= 3 {
			break
		}
	}

	cardList = append(cardList, model.TextCard{
		Title:       "4호선(한대앞역)",
		Description: strings.TrimSpace(message),
		Buttons:     []model.CardButton{},
	})

	message = "수원 방면\n"
	for i, item := range lineSuinBundang.Realtime.Up {
		message += item.Terminal + "행 " + strconv.Itoa(item.Time) + "분 후 도착 (" + item.Location + ")\n"
		if i >= 3 {
			break
		}
	}
	message += "\n인천 방면\n"
	for i, item := range lineSuinBundang.Realtime.Down {
		message += item.Terminal + "행 " + strconv.Itoa(item.Time) + "분 후 도착 (" + item.Location + ")\n"
		if i >= 3 {
			break
		}
	}

	cardList = append(cardList, model.TextCard{
		Title:       "수인분당선(한대앞역)",
		Description: strings.TrimSpace(message),
		Buttons:     []model.CardButton{},
	})
	response := util.SetResponse(util.SetTemplate([]model.Components{
		util.SetBasicCardCarousel(cardList)}, []model.QuickReply{}))
	return c.JSON(response)
}
