//+build server

package api

import (
	"OpenGreenEye/database"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type ApiDriver struct {
	*fiber.App
	JSONDencoder func(data []byte, ptr interface{}) error
	Validate     func(Struct interface{}) error
}

type routerDriver struct {
	fiber.Router
	JSONDencoder func(data []byte, ptr interface{}) error
	Validate     func(Struct interface{}) error
	DB           database.DriverGlobal
}

type routesFunc func(router routerDriver)

var publicRoutes = make([]routesFunc, 0)

var sensorsRoutes = make([]routesFunc, 0)

var clientRoutes = make([]routesFunc, 0)

const KEY_HEADER_NAME = "KEY"

func LoadRoutes(api *ApiDriver, db database.DriverGlobal) {
	publicRouter := routerDriver{
		Router:       nil,
		JSONDencoder: api.JSONDencoder,
		Validate:     api.Validate,
		DB:           db,
	}

	publicRouter.Router = api.Group("/api/v1/public")

	for _, f := range publicRoutes {
		f(publicRouter)
	}

	sensorsRouter := routerDriver{
		Router:       nil,
		JSONDencoder: api.JSONDencoder,
		DB:           db,
	}

	sensorsRouter.Router = api.Group("/api/v1/sensors", func(ctx *fiber.Ctx) error {
		if key := ctx.Get(KEY_HEADER_NAME, ""); key != "" {
			if id := sensorsRouter.DB.CheckSensorKey(key, ctx.IP()); id >= 0 {
				ctx.Request().Header.Set(KEY_HEADER_NAME, strconv.FormatInt(int64(id), 10))
				return ctx.Next()
			}
		}

		return ctx.SendStatus(http.StatusNotFound)
	})

	for _, f := range sensorsRoutes {
		f(sensorsRouter)
	}

	clientRouter := routerDriver{
		Router:       nil,
		JSONDencoder: api.JSONDencoder,
		DB:           db,
	}

	clientRouter.Router = api.Group("/api/v1/client", func(ctx *fiber.Ctx) error {
		if key := ctx.Get(KEY_HEADER_NAME, ""); key != "" {
			if clientRouter.DB.CheckUserKey(key) {
				return ctx.Next()
			}
		}

		return ctx.SendStatus(http.StatusNotFound)
	})

	for _, f := range clientRoutes {
		f(clientRouter)
	}
}

func getJsonBody(r *routerDriver, ctx *fiber.Ctx, ptr interface{}) error {
	if err := r.JSONDencoder(ctx.Body(), ptr); err != nil {
		return err
	}

	if err := r.Validate(ptr); err != nil {
		return err
	}

	return nil
}
