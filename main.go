package main

import (
	"Server/handler"
	"Server/models"
	"github.com/labstack/echo/v4"
	"sort"
)

type User []models.GithubUserInfo

func (a User) Len() int           { return len(a) }
func (a User) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a User) Less(i, j int) bool { return a[i].Login < a[j].Login }

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
		followingNum, followersNum := handler.GetUserFollowInfo(userName)
		followingList = User(handler.GetFollowUserList(userName, "following", followingNum))
		followersList = User(handler.GetFollowUserList(userName, "followers", followersNum))
		for _, user := range followingList {
			handler.UserSet1(user, m)
		}
		for _, user := range followersList {
			a := handler.FindUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(list)
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
		followingNum, followersNum := handler.GetUserFollowInfo(userName)
		followingList = User(handler.GetFollowUserList(userName, "following", followingNum))
		followersList = User(handler.GetFollowUserList(userName, "followers", followersNum))
		for _, user := range followersList {
			handler.UserSet1(user, m)
		}
		for _, user := range followingList {
			a := handler.FindUnfollwer(user, m)
			if a == 1 {
				list = append(list, user)
			}
		}
		sort.Sort(list)
		return c.JSON(200, list)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
