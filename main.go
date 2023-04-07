package main

import (
	"Server/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type findUser struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}

type ByLogin []User

func (a ByLogin) Len() int           { return len(a) }
func (a ByLogin) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLogin) Less(i, j int) bool { return a[i].Login < a[j].Login }

var baseurl = "https://api.github.com/users/"

func main() {

	e := echo.New()
	e.GET("/unfollowing", func(c echo.Context) error {
		userName := c.QueryParam("userName")
		//unfollowCh := make(chan User)
		m := make(map[string]int)
		var (
			followingList []User
			followersList []User
			list          []User
		)
		followingNum, followersNum := getUserFollowInfo(userName)
		followingList = getFollowUserList(userName, "following", followingNum)
		followersList = getFollowUserList(userName, "followers", followersNum)
		for _, user := range followingList {
			userSet1(user, m)
		}
		for _, user := range followersList {
			a := findUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(ByLogin(list))
		return c.JSON(200, list)
	})

	e.GET("/unfollower", func(c echo.Context) error {
		userName := c.QueryParam("userName")
		//unfollowCh := make(chan User)
		m := make(map[string]int)
		var (
			followingList []User
			followersList []User
			list          []User
		)
		followingNum, followersNum := getUserFollowInfo(userName)
		followingList = getFollowUserList(userName, "following", followingNum)
		followersList = getFollowUserList(userName, "followers", followersNum)
		for _, user := range followersList {
			userSet1(user, m)
		}
		for _, user := range followingList {
			a := findUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(ByLogin(list))
		return c.JSON(200, list)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

// 팔로잉, 팔로워 두 함수를 하나로 합침
func getFollowUserList(userName string, follow string, length int) []User {
	userLen := length
	if length%100 != 0 {
		length = length/100 + 1
	} else {
		length = length / 100
	}
	list := make([]User, 0)
	c := make(chan User)
	for i := 1; i <= length; i++ {
		go hitURL(userName, follow, i, c)
	}
	for i := 0; i < userLen; i++ {
		user := <-c
		list = append(list, user)
	}
	return list
}

// 내 팔로잉 팔로워 갯수를 알아오는 함수
func getUserFollowInfo(userName string) (int, int) {
	url := baseurl + userName
	req, err := http.NewRequest("GET", url, nil)
	utils.CheckErr(err)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}

	res, err := client.Do(req)
	utils.CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	var user findUser
	err = json.Unmarshal(body, &user)
	utils.CheckErr(err)
	return user.Following, user.Followers
}

// 맵에 유저의 정보를 담아줄 함수
// 문제 없는거 확인
func userSet1(user User, m map[string]int) {
	m[user.Login] = 1
}

// 맵에 user가 들어있는지 확인해줄 함수
func findUnfollwer(user User, m map[string]int) int {
	val := m[user.Login]
	if val != 1 {
		return 1
	}
	return 0
}

func hitURL(userName string, follow string, i int, c chan User) {
	pageURL := baseurl + userName + "/" + follow + "?per_page=100&page=" + strconv.Itoa(i)
	req, err := http.NewRequest("GET", pageURL, nil)
	utils.CheckErr(err)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	res, err := client.Do(req)
	utils.CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	var following []User
	err = json.Unmarshal(body, &following)
	utils.CheckErr(err)
	for _, user := range following {
		c <- user
	}
}
