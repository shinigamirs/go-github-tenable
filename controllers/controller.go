package controllers

import (
	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type ParamValidator struct {
	validator *validator.Validate
}

func (pv *ParamValidator) Validate(i interface{}) error {
	if err := pv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func InitEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &ParamValidator{validator: validator.New()}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	return e
}

func AddRoutes(e *echo.Echo) {
	AddGithubRoute(e)
	AddLoginRoute(e)
}
