package main

import (
	"github.com/pscheid92/idcard-generator/internal/middleware"
	"github.com/pscheid92/idcard-generator/internal/models"
	"github.com/pscheid92/idcard-generator/internal/renderer"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Renderer = renderer.NewTemplateRenderer("templates/*")
	e.Use(middleware.ForwardedPrefixMiddleware)

	// main page
	e.GET("/", getMainPage)
	e.POST("/", postMainPage)

	// health endpoint
	e.GET("/health", healthEndpoint)

	err := e.Start(":8080")
	e.Logger.Fatal(err)
}

func healthEndpoint(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func getMainPage(c echo.Context) error {
	model := models.NewViewModel()
	model.CardOptions[0].Selected = true
	model.PathPrefix = c.Get("PathPrefix").(string)
	return c.Render(http.StatusOK, "index.html", model)
}

func postMainPage(c echo.Context) error {
	form := new(postModel)
	if err := c.Bind(form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	model := models.NewViewModel()
	model.PathPrefix = c.Get("PathPrefix").(string)
	model.Birthday = form.Birthday
	model.Expiration = form.Expiration
	model.Manipulation = form.Manipulation == "manipulation"

	switch form.CardType {
	case models.NewId:
		model.CardOptions[0].Selected = true
		model.CalculateNewId()
	case models.OldId:
		model.CardOptions[1].Selected = true
		model.CalculateOldId()
	case models.Passport:
		model.CardOptions[2].Selected = true
		model.CalculatePassport()
	}

	return c.Render(http.StatusOK, "index.html", model)
}

type postModel struct {
	Birthday     string `form:"birthday"`
	Expiration   string `form:"expiration"`
	Manipulation string `form:"manipulation"`
	CardType     string `form:"cardtype"`
}
