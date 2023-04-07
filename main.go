package main

import (
	"Server/handler"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.GET("/unfollowing", handler.UnfollowingCheckFunc)

	e.GET("/unfollower", handler.UnfollowerCheckFunc)
	e.Logger.Fatal(e.Start(":8080"))
}
