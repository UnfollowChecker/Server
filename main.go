package main

import (
	"Server/utils"
	"encoding/json"
	"fmt"
	"net/http"
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

//var baseurl string = "https://api.github.com/users/yoochanhong/following?per_page=100"

func main() {
	//e := echo.New()
	//e.Logger.Fatal(e.Start(":8080"))
	//var res *http.Response
	//for i := 1; i < 3; i++ {
	//pageURL := baseurl + "&page=" + strconv.Itoa(i)
	res, err := http.Get("https://api.github.com/users/yoochanhong/following?per_page=100")
	res.Header.Set("Authorization", "Bearer"+token)
	utils.CheckErr(err)
	var users []User
	err = json.NewDecoder(res.Body).Decode(&users)
	utils.CheckErr(err)
	for _, user := range users {
		fmt.Println(user.Login)
	}
	//}
}
