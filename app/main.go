package main

import (
	_ "database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/warlikegh/db_project/config"
)

type initRoute struct {
	e *echo.Echo

}

func main() {
	fmt.Println(config.Host)
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339_nano}\n",
	}))

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
