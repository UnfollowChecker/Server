package main

import "github.com/labstack/echo/v4"

type User struct {
	login             string
	id                int
	nodeId            string
	avatarUrl         string
	gravatarId        string
	url               string
	htmlUrl           string
	followersUrl      string
	followingUrl      string
	gistsUrl          string
	starredUrl        string
	subscriptionsUrl  string
	organizationsUrl  string
	reposUrl          string
	eventsUrl         string
	receivedEventsUrl string
	userType          string
	siteAdmin         string
}

func main() {
	e := echo.New()
	e.Logger.Fatal(e.Start(":8080"))
}
