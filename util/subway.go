package util

import (
	"encoding/json"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"
	"io"
	"net/http"
	"os"
)

func GetSubwayArrival(stationID string) model.SubwayStationResponse {
	response, err := http.Get(os.Getenv("API_URL") + "/rest/subway/station/" + stationID)
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
	arrivalResponse := model.SubwayStationResponse{}
	err = json.Unmarshal(data, &arrivalResponse)
	if err != nil {
		panic(err)
	}
	return arrivalResponse
}
