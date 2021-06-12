package main

import (
	"database/sql"
	"fmt"
	"forums/config"
	"forums/internal/forum"
	forumDelivery "forums/internal/forum/delivery"
	forumRepo "forums/internal/forum/repository"
	forumUcase "forums/internal/forum/usecase"
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
	e     *echo.Echo
	user  user.UserHandler
	forum forum.ForumHandler
}

func handler(c echo.Context) error {
	fmt.Println("paramValues = ", c.ParamValues(), "\nPath = ", c.Path())
	return nil
}

func route(data initRoute) {
	data.e.POST("/api/forum/create", data.forum.CreateForum)
	data.e.GET("/api/forum/:slug/details", data.forum.GetForumBySlug)
	data.e.POST("/api/forum/:slug/create", data.forum.CreateThread)
	data.e.GET("/api/forum/:slug/users", handler)
	data.e.GET("/api/forum/:slug/threads", handler)

	data.e.GET("/api/post/:id/details", handler)
	data.e.POST("/api/post/:id/details", handler)

	data.e.POST("/api/service/clear", handler)
	data.e.GET("/api/service/status", handler)

	data.e.POST("/api/thread/:slugOrId/create", handler)
	data.e.GET("/api/thread/:slugOrId/details", handler)
	data.e.POST("/api/thread/:slugOrId/details", handler)
	data.e.GET("/api/thread/:slugOrId/posts", handler)
	data.e.POST("/api/thread/:slugOrId/vote", handler)

	data.e.POST("/api/user/:nickname/create", data.user.Create)
	data.e.GET("/api/user/:nickname/profile", data.user.GetUserData)
	data.e.POST("/api/user/:nickname/profile", data.user.UpdateUserData)
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

	forumRepo_ := forumRepo.NewForumRepo(db, userRepo_)
	forumUcase_ := forumUcase.NewForumUsecase(forumRepo_)
	forumHandler_ := forumDelivery.NewForumHandler(forumUcase_)

	route(initRoute{
		e:     e,
		user:  userHandler_,
		forum: forumHandler_,
	})

	e.Logger.Fatal(e.Start(":5000"))
}
