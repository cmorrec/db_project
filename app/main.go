package main

import (
	"database/sql"
	"fmt"
	"forums/config"
	"forums/internal/user"
	userDelivery "forums/internal/user/delivery/http"
	userRepo "forums/internal/user/repository"
	userUcase "forums/internal/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
)

type initRoute struct {
	e    *echo.Echo
	user user.UserHandler
}

func handler(c echo.Context) error {
	fmt.Println("paramValues = ", c.ParamValues(), "\nPath = ", c.Path())
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

	data.e.POST("/user/:nickname/create", data.user.Create)
	data.e.GET("/user/:nickname/profile", data.user.GetUserData)
	data.e.POST("/user/:nickname/profile", handler)
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}\n uri=${uri}\n status=${status}\n time=${time_rfc3339_nano}\n\n",
	}))

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRepo_ := userRepo.NewUserRepo(db)
	userUcase_ := userUcase.NewUserUsecase(userRepo_)
	userHandler_ := userDelivery.NewUserHandler(userUcase_)

	route(initRoute{
		e:    e,
		user: userHandler_,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
