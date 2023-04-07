package handler

import (
	"Server/models"
	"Server/private"
	"Server/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

var baseurl = "https://api.github.com/users/"

type User []models.GithubUserInfo

// 맵에 유저의 정보를 담아줄 함수
func UserSet1(user models.GithubUserInfo, m map[string]int) {
	m[user.Login] = 1
}

// 맵에 user가 들어있는지 확인해줄 함수
func FindUnfollwer(user models.GithubUserInfo, m map[string]int) int {
	val := m[user.Login]
	if val != 1 {
		return 1
	}
	return 0
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
func GetFollowUserList(userName string, follow string, length int) User {
	userLen := length
	if length%100 != 0 {
		length = length/100 + 1
	} else {
		length = length / 100
	}
	list := make(User, 0)
	c := make(chan models.GithubUserInfo)
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
func GetUserFollowInfo(userName string) (int, int) {
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
