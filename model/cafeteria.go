package model

type RestaurantListResponse struct {
	Restaurant []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Menu []struct {
			Time     string `json:"time"`
			MenuList []struct {
				Food  string `json:"food"`
				Price string `json:"price"`
			} `json:"menu"`
		} `json:"menu"`
	} `json:"restaurant"`
}
