package main

import (
	"Server/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
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

type findUser struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}

var baseurl = "https://api.github.com/users/"

func main() {

	e := echo.New()

	e.GET("/unfollower", func(c echo.Context) error {
		userName := c.QueryParam("userName")
		followingNum, followersNum := getUserFollowInfo(userName)
		fmt.Println(followersNum, followingNum)
		return c.JSON(200, "what")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

// 팔로잉, 팔로워 두 함수를 하나로 합침
func getFollowUserList(userName string, follow string, length int) []User {
	var list []User
	for i := 1; length > 0; i++ {
		pageURL := baseurl + userName + "/" + follow + "?per_page=100&page=" + strconv.Itoa(i)
		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		var following []User
		err = json.Unmarshal(body, &following)
		utils.CheckErr(err)
		for _, user := range following {
			list = append(list, user)
		}
		length -= 100
	}
	return list
}

// 내 팔로잉 팔로워 갯수를 알아오는 함수
func getUserFollowInfo(userName string) (int, int) {
	url := baseurl + userName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	var user findUser
	err = json.Unmarshal(body, &user)
	utils.CheckErr(err)
	return user.Following, user.Followers
}

// 맵에 유저의 정보를 담아줄 함수
// 문제 없는거 확인
func userSet1(user User, m map[string]int, mutex *sync.Mutex) {
	//mutex.Lock()
	m[user.Login] = 1
	//mutex.Unlock()
}

// 맵에 user가 들어있는지 확인해줄 함수
func findUnfollwer(user User, m map[string]int, mutex *sync.Mutex, ch chan User) {
	mutex.Lock()
	val := m[user.Login]
	if val != 1 {
		ch <- user
	}
	mutex.Unlock()
}
