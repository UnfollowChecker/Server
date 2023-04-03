package main

import (
	"Server/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

var baseurl = "https://api.github.com/users/"

func main() {

	var (
		followingList []User
		followerList  []User
	)

	m := make(map[string]int)
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "what")
	})

	e.GET("/unfollower", func(c echo.Context) error {

		userName := c.QueryParam("userName")
		followingList = getFollowUserList(userName)
		followerList = getFollowerUserList(userName)

		var list []User
		for _, user := range followerList {
			m[user.Login] = 1
		}
		for _, user := range followingList {
			if m[user.Login] != 1 {
				list = append(list, user)
			}
		}
		return c.JSON(200, list)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

// 내가 팔로잉한 사람들을 긁어오는 함수
func getFollowUserList(userName string) []User {
	var userList []User
	for i := 1; ; i++ {
		pageURL := baseurl + userName + "/following?per_page=100&page=" + strconv.Itoa(i)
		req, err := http.NewRequest("GET", pageURL, nil)
		utils.CheckErr(err)
		req.Header.Set("Authorization", "Bearer"+token)
		client := &http.Client{}

		res, err := client.Do(req)
		utils.CheckErr(err)
		var users []User
		err = json.NewDecoder(res.Body).Decode(&users)
		utils.CheckErr(err)
		for _, user := range users {
			userList = append(userList, user)
		}
		if len(users) != 100 {
			break
		}
	}
	return userList
}

// 내 팔로워를 모두 가져오는 함수
func getFollowerUserList(userName string) []User {
	var userList []User
	for i := 1; ; i++ {
		pageURL := baseurl + userName + "/followers?per_page=100&page=" + strconv.Itoa(i)
		req, err := http.NewRequest("GET", pageURL, nil)
		utils.CheckErr(err)
		req.Header.Set("Authorization", "Bearer"+token)
		client := &http.Client{}

		res, err := client.Do(req)
		utils.CheckErr(err)
		var users []User
		err = json.NewDecoder(res.Body).Decode(&users)
		utils.CheckErr(err)
		for _, user := range users {
			userList = append(userList, user)
		}
		if len(users) != 100 {
			break
		}
	}
	return userList
}
