package main

import (
	"Server/utils"
	"encoding/json"
	"fmt"
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

var baseurl string = "https://api.github.com/users/"

func main() {
	//e := echo.New()
	//e.Logger.Fatal(e.Start(":8080"))
	var (
		followingList []User
		//unFollowerList []User
	)
	followingList = getFollowUserList("yoochanhong")
	//unFollowerList = getUnFollowUserList("yoochanhong", followingList)
	for _, user := range followingList {
		fmt.Println(user.Login)
	}
	//for _, user := range unFollowerList {
	//	fmt.Println(user.Login)
	//}
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

// 내가 팔로우했지만 나를 팔로우하지 않은 사람들을 모두 가져오는 함수
//func getUnFollowUserList(userName string, followList []User) []User {
//	var userList []User
//	for _, user := range followList {
//		pageURL := baseurl + user.Login + "/following/" + userName
//		res, err := http.Get(pageURL)
//		res.Header.Set("Authorization", "Bearer"+token)
//		utils.CheckErr(err)
//		if res.StatusCode == 404 {
//			userList = append(userList, user)
//		}
//	}
//	return userList
//}
