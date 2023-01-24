package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/util"
)

func CafeteriaMenu(c *fiber.Ctx) error {
	restaurantList := util.GetCafeteriaMenu(2)
	cardList := make([]model.TextCard, 0)
	quickReplies := make([]model.QuickReply, 0)
	for _, restaurant := range restaurantList.Restaurant {
		description := ""
		if len(restaurant.Menu) == 0 {
			description = "오늘은 메뉴가 없어요!"
		} else {
			for _, menu := range restaurant.Menu[0].MenuList {
				description += menu.Food + "\n"
				description += menu.Price + "원\n"
			}
		}
		cardList = append(cardList, model.TextCard{
			Title:       restaurant.Name,
			Description: description,
			Buttons:     []model.CardButton{},
		})
	}
	response := util.SetResponse(
		util.SetTemplate([]model.Components{util.SetBasicCardCarousel(cardList)}, quickReplies))
	return c.JSON(response)
}
