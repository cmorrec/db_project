package main

import (
	_ "database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/warlikegh/db_project/config"
)

type initRoute struct {
	e *echo.Echo
}

func handler(c echo.Context) error {
	fmt.Println("paramValues = ", c.ParamValues(), "\nPath = ", c.Path(), "\nqueryParams = ", c.QueryParams(), "\nRequest = ",
		c.Request(), "\n\n")
	return nil
}

func route(data initRoute) {
	data.e.POST("/forum/create", handler)
	data.e.GET("/forum/:slug/details", handler)
	data.e.POST("/forum/:slug/create", handler)
	data.e.GET("/forum//:slug/users", handler)
	data.e.GET("/forum/:slug/threads", handler)

	data.e.GET("/post/:id/details", handler)
	data.e.POST("/post/:id/details", handler)

	data.e.POST("/service/clear", handler)
	data.e.GET("/service/status", handler)

	data.e.POST("/thread/:slugOrId/create", handler)
	data.e.GET("/thread/:slugOrId/details", handler)
	data.e.POST("/thread/:slugOrId/details", handler)
	data.e.GET("/thread/:slugOrId/posts", handler)
	data.e.POST("/thread/:slugOrId/vote", handler)

	data.e.POST("/user/:nickname/create", handler)
	data.e.GET("/user/:nickname/profile", handler)
	data.e.POST("/user/:nickname/profile", handler)
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339_nano}\n",
	}))
	route(initRoute{e})

	//dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	//db, err := sql.Open(config.PostgresDB, dsn)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(3)
	//
	//err = db.Ping()
	//if err != nil {
	//	log.Fatal(err)
	//}

	e.Logger.Fatal(e.Start(":5000"))
}
