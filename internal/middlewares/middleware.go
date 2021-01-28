package middlewares

import "github.com/labstack/echo/v4"

func ForwardedPrefixMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		pathPrefix := c.Request().Header.Get("X-Forwarded-Prefix")
		if pathPrefix == "" {
			pathPrefix = "/"
		}

		c.Set("PathPrefix", pathPrefix)

		return next(c)
	}
}
