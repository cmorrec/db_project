package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
}

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

func SendResponseWithErrorConflictMessage(c echo.Context) error {
	return c.JSON(http.StatusConflict, CustomError{"Can't find user with id #42\n"})
}
