package handler

import (
	"Server/models"
	"Server/private"
	"Server/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type User []models.GithubUserInfo

func (a User) Len() int           { return len(a) }
func (a User) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a User) Less(i, j int) bool { return a[i].Login < a[j].Login }

var baseurl = "https://api.github.com/users/"

// 맵에 유저의 정보를 담아줄 함수
func userSet1(user models.GithubUserInfo, m map[string]int) {
	m[user.Login] = 1
}

// 맵에 user가 들어있는지 확인해줄 함수
func findUnfollwer(user models.GithubUserInfo, m map[string]int, list *User) {
	val := m[user.Login]
	if val != 1 {
		*list = append(*list, user)
	}
}

func hitURL(userName string, follow string, i int, c chan models.GithubUserInfo) {
	pageURL := baseurl + userName + "/" + follow + "?per_page=100&page=" + strconv.Itoa(i)
	req, err := http.NewRequest("GET", pageURL, nil)
	utils.CheckErr(err)
	req.Header.Set("Authorization", "Bearer "+private.Token)
	client := &http.Client{}
	res, err := client.Do(req)
	utils.CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	var following User
	err = json.Unmarshal(body, &following)
	utils.CheckErr(err)
	for _, user := range following {
		c <- user
	}
}

// 팔로잉, 팔로워 두 함수를 하나로 합침
func getFollowUserList(userName string, follow string, length int, list *User, ch chan models.GithubUserInfo, ch1 chan int) {
	userLen := length
	if length%100 != 0 {
		length = length/100 + 1
	} else {
		length = length / 100
	}
	for i := 1; i <= length; i++ {
		go hitURL(userName, follow, i, ch)
	}
	for i := 0; i < userLen; i++ {
		*list = append(*list, <-ch)
	}
	ch1 <- 1
}

// 내 팔로잉 팔로워 갯수를 알아오는 함수
func getUserFollowInfo(userName string) (int, int) {
	url := baseurl + userName
	req, err := http.NewRequest("GET", url, nil)
	utils.CheckErr(err)
	req.Header.Set("Authorization", "Bearer "+private.Token)
	client := &http.Client{}

	res, err := client.Do(req)
	utils.CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	var user models.FindUser
	err = json.Unmarshal(body, &user)
	utils.CheckErr(err)
	return user.Following, user.Followers
}

func UnfollowingCheckFunc(c echo.Context) error {
	userName := c.QueryParam("userName")
	followingCh := make(chan models.GithubUserInfo)
	followersCh := make(chan models.GithubUserInfo)
	ch := make(chan int)
	m := make(map[string]int)
	var (
		followingList User
		followersList User
		list          User
	)

	followingNum, followersNum := getUserFollowInfo(userName)
	fmt.Println(followingNum, followersNum)
	go getFollowUserList(userName, "following", followingNum, &followingList, followingCh, ch)
	go getFollowUserList(userName, "followers", followersNum, &followersList, followersCh, ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	for _, user := range followingList {
		userSet1(user, m)
	}
	for _, user := range followersList {
		findUnfollwer(user, m, &list)
	}
	sort.Sort(list)
	return c.JSON(200, list)
}

func UnfollowersCheckFunc(c echo.Context) error {
	userName := c.QueryParam("userName")
	followingCh := make(chan models.GithubUserInfo)
	followersCh := make(chan models.GithubUserInfo)
	ch := make(chan int)
	m := make(map[string]int)
	var (
		followingList User
		followersList User
		list          User
	)

	followingNum, followersNum := getUserFollowInfo(userName)
	go getFollowUserList(userName, "following", followingNum, &followingList, followingCh, ch)
	go getFollowUserList(userName, "followers", followersNum, &followersList, followersCh, ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	for _, user := range followersList {
		userSet1(user, m)
	}
	for _, user := range followingList {
		findUnfollwer(user, m, &list)
	}
	sort.Sort(list)
	return c.JSON(200, list)
}
