package main

import (
	"Server/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type User struct {
	Login             string
	Id                int
	NodeId            string
	AvatarUrl         string
	GravatarId        string
	Url               string
	HtmlUrl           string
	FollowersUrl      string
	FollowingUrl      string
	GistsUrl          string
	StarredUrl        string
	SubscriptionsUrl  string
	OrganizationsUrl  string
	ReposUrl          string
	EventsUrl         string
	ReceivedEventsUrl string
	Type              string
	SiteAdmin         bool
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
		followingList = getFollowUserList("yoochanhong")
		followerList = getFollowerUserList("yoochanhong")

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
