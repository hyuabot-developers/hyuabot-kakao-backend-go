package util

import (
	"encoding/json"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetCafeteriaMenu(campusID int) model.RestaurantListResponse {
	response, err := http.Get(os.Getenv("API_URL") + "/rest/cafeteria/campus/" + strconv.Itoa(campusID) + "/restaurant")
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	restaurantData := model.RestaurantListResponse{}
	err = json.Unmarshal(data, &restaurantData)
	if err != nil {
		panic(err)
	}
	return restaurantData
}
