package model

type StopItemResponse struct {
	Id           int          `json:"id"`
	Name         string       `json:"name"`
	MobileNumber string       `json:"mobileNumber"`
	Location     StopLocation `json:"location"`
	Route        []StopRoute  `json:"route"`
}

type StopLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type StopRoute struct {
	Id       int                        `json:"id"`
	Name     string                     `json:"name"`
	Start    StopRouteStartStop         `json:"start"`
	Realtime []StopRouteRealtimeArrival `json:"realtime"`
}

type StopRouteStartStop struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Timetable []string `json:"timetable"`
}

type StopRouteRealtimeArrival struct {
	Stop     int  `json:"stop"`
	Time     int  `json:"time"`
	Seat     int  `json:"seat"`
	LowPlate bool `json:"lowPlate"`
}
