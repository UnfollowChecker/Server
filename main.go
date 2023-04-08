package main

import (
	"Server/handler"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.GET("/unfollowing", handler.UnfollowingCheckFunc)

	e.GET("/unfollowers", handler.UnfollowersCheckFunc)
	e.Logger.Fatal(e.Start(":8080"))
}
