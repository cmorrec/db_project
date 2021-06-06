package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
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
	serverMessage := message{http.StatusOK, data}
	return c.JSON(http.StatusOK, serverMessage)
}

func SendResponseWithError(c echo.Context, err error) error {
	// TODO
	return err
}

