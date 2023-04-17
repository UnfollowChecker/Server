package models

type UserFollowInfo struct {
	Data struct {
		User struct {
			Followers struct {
				TotalCount int `json:"totalCount"`
			} `json:"followers"`
			Following struct {
				TotalCount int `json:"totalCount"`
			} `json:"following"`
		} `json:"user"`
	} `json:"data"`
}
