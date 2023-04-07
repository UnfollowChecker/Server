package main

import (
	"Server/handler"
	"Server/models"
	"Server/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sort"
)

type User []models.GithubUserInfo

func (a User) Len() int           { return len(a) }
func (a User) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a User) Less(i, j int) bool { return a[i].Login < a[j].Login }

var baseurl = "https://api.github.com/users/"

func main() {

	e := echo.New()
	e.GET("/unfollowing", func(c echo.Context) error {
		userName := c.QueryParam("userName")
		//unfollowCh := make(chan User)
		m := make(map[string]int)
		var (
			followingList User
			followersList User
			list          User
		)
		followingNum, followersNum := getUserFollowInfo(userName)
		followingList = getFollowUserList(userName, "following", followingNum)
		followersList = getFollowUserList(userName, "followers", followersNum)
		for _, user := range followingList {
			handler.UserSet1(user, m)
		}
		for _, user := range followersList {
			a := handler.FindUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(User(list))
		return c.JSON(200, list)
	})

	e.GET("/unfollower", func(c echo.Context) error {
		userName := c.QueryParam("userName")
		//unfollowCh := make(chan User)
		m := make(map[string]int)
		var (
			followingList User
			followersList User
			list          User
		)
		followingNum, followersNum := getUserFollowInfo(userName)
		followingList = getFollowUserList(userName, "following", followingNum)
		followersList = getFollowUserList(userName, "followers", followersNum)
		for _, user := range followersList {
			handler.UserSet1(user, m)
		}
		for _, user := range followingList {
			a := handler.FindUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(User(list))
		return c.JSON(200, list)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

// 팔로잉, 팔로워 두 함수를 하나로 합침
func getFollowUserList(userName string, follow string, length int) User {
	userLen := length
	if length%100 != 0 {
		length = length/100 + 1
	} else {
		length = length / 100
	}
	list := make(User, 0)
	c := make(chan models.GithubUserInfo)
	for i := 1; i <= length; i++ {
		go handler.HitURL(userName, follow, i, c)
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
	var user models.FindUser
	err = json.Unmarshal(body, &user)
	utils.CheckErr(err)
	return user.Following, user.Followers
}
