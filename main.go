package main

import (
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pscheid92/idcard-generator/internal/middlewares"
	"github.com/pscheid92/idcard-generator/internal/models"
	"github.com/pscheid92/idcard-generator/internal/renderer"
)

func main() {
	secureMiddlewareConfig := middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            315360000,
		HSTSExcludeSubdomains: true,
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = renderer.NewTemplateRenderer("templates/*")
	e.Use(middlewares.ForwardedPrefixMiddleware)
	e.Use(middleware.SecureWithConfig(secureMiddlewareConfig))

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
	case models.NewID:
		model.CardOptions[0].Selected = true
		model.CalculateNewID()
	case models.OldID:
		model.CardOptions[1].Selected = true
		model.CalculateOldID()
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
