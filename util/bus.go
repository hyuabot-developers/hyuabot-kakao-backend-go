package util

import (
	"encoding/json"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetBusArrival(busStopID int) model.StopItemResponse {
	response, err := http.Get(os.Getenv("API_URL") + "/rest/bus/stop/" + strconv.Itoa(busStopID))
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
	busStopData := model.StopItemResponse{}
	err = json.Unmarshal(data, &busStopData)
	if err != nil {
		panic(err)
	}
	return busStopData
}
