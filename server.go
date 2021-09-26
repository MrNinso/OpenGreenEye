//+build server

package main

import (
	"OpenGreenEye/api"
	"OpenGreenEye/database"
	"OpenGreenEye/validate"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"os"
)

func main() {
	db := database.NewClient()

	j := jsoniter.Config{
		MarshalFloatWith6Digits: false,
		SortMapKeys:             false,
		UseNumber:               false,
		DisallowUnknownFields:   false,
		OnlyTaggedField:         true,
		CaseSensitive:           true,
	}.Froze()

	app := fiber.New(fiber.Config{
		DisableKeepalive:         true,
		DisableHeaderNormalizing: true,
		DisableStartupMessage:    true,
		JSONEncoder:              j.Marshal,
	})

	apiDriver := api.ApiDriver{
		App:          app,
		JSONDencoder: j.Unmarshal,
		Validate:     validate.ConfigValidate().Struct,
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	api.LoadRoutes(&apiDriver, db)

	if err := apiDriver.Listen(fmt.Sprint(":", port)); err != nil {
		panic(err)
	}
}
