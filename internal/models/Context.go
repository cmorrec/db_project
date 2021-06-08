package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
}
//
//func GetContext(c echo.Context) context.Context {
//	ctx := c.Request().Context()
//
//	ctx = context.WithValue(ctx, "User", c.Get("User"))
//	ctx = context.WithValue(ctx, "Restaurant", c.Get("Restaurant"))
//	return context.WithValue(ctx, "request_id", c.Get("request_id"))
//}



func SendResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func SendResponseCreate(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}

func SendResponseWithErrorNotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, CustomError{"Can't find user with id #42\n"})
}

func SendResponseWithErrorConflict(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusConflict, data)
}

