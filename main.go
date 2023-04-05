package main

import (
	"Server/utils"
	"encoding/json"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "what")
	})

	e.GET("/unfollower", func(c echo.Context) error {
		m := make(map[int]int)
		unfollowerCh := make(chan User)
		var (
			followingList []User
			followerList  []User
			list          []User
		)

		userName := c.QueryParam("userName")
		followingNum, followerNum := getUserFollowInfo(userName)
		followingList = getFollowUserList(userName, "following", followingNum)
		followerList = getFollowUserList(userName, "followers", followerNum)
		mutex := sync.Mutex{}

		for _, user := range followerList {
			go userSet1(user, m, &mutex)
		}
		for _, user := range followingList {
			go findUnfollwer(user, m, &mutex, unfollowerCh)
			go func() {
				list = append(list, <-unfollowerCh)
			}()
		}
		return c.JSON(200, list)
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
func userSet1(user User, m map[int]int, mutex *sync.Mutex) {
	mutex.Lock()
	m[user.ID] = 1
	mutex.Unlock()
}

// 맵에 user가 들어있는지 확인해줄 함수
func findUnfollwer(user User, m map[int]int, mutex *sync.Mutex, ch chan User) {
	mutex.Lock()
	if m[user.ID] != 1 {
		ch <- user
	}
	mutex.Unlock()
}
