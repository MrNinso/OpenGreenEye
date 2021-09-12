//+build server auth

package api

import (
	"OpenGreenEye/database"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func init() {
	publicRoutes = append(publicRoutes, BuildAuthRoute)
}

func BuildAuthRoute(driver routerDriver) {

	DB := driver.DB.(database.DriverAuth)

	driver.Router.Post("/login", func(ctx *fiber.Ctx) error {
		var b struct {
			Login    string `json:"login" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := getJsonBody(&driver, ctx, &b); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := DB.CheckLoginUser(b.Login, b.Password)

		if token == "" {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		return ctx.JSON(fiber.Map{
			"token": token,
		})
	})

	driver.Router.Post("/register", func(ctx *fiber.Ctx) error {
		var b struct {
			Login    string `json:"login" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := getJsonBody(&driver, ctx, &b); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token, err := DB.RegisterUser(b.Login, b.Password)

		if token == "" || err != uint8(0) {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"token": "",
				"error": err,
			})
		}

		return ctx.JSON(fiber.Map{
			"token": token,
			"error": 0,
		})
	})
}
