package model

type ShuttleArrivalResponse struct {
	Stop []ShuttleArrivalStopResponse `json:"stop"`
}

type ShuttleArrivalStopResponse struct {
	Name  string                            `json:"name"`
	Route []ShuttleArrivalStopRouteResponse `json:"route"`
}

type ShuttleArrivalStopRouteResponse struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Arrival []int  `json:"arrival"`
}
