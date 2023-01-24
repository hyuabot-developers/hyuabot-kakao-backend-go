package model

type SubwayStationResponse struct {
	StationID   string `json:"stationID"`
	StationName string `json:"stationName"`
	RouteID     int    `json:"routeID"`
	Realtime    struct {
		Up []struct {
			Terminal    string `json:"terminal"`
			TrainNumber string `json:"trainNumber"`
			Time        int    `json:"time"`
			Stop        int    `json:"stop"`
			Express     bool   `json:"express"`
			Heading     bool   `json:"heading"`
			Location    string `json:"location"`
		} `json:"up"`
		Down []struct {
			Terminal    string `json:"terminal"`
			TrainNumber string `json:"trainNumber"`
			Time        int    `json:"time"`
			Stop        int    `json:"stop"`
			Express     bool   `json:"express"`
			Heading     bool   `json:"heading"`
			Location    string `json:"location"`
		} `json:"down"`
	} `json:"realtime"`
	Timetable struct {
		Up []struct {
			Start         string `json:"start"`
			Terminal      string `json:"terminal"`
			DepartureTime string `json:"departureTime"`
		} `json:"up"`
		Down []struct {
			Start         string `json:"start"`
			Terminal      string `json:"terminal"`
			DepartureTime string `json:"departureTime"`
		} `json:"down"`
	} `json:"timetable"`
}
