package routes

import (
	"authentication/auth"
	"authentication/base"

	"github.com/gofiber/fiber/v2"
)

func Install(app *fiber.App) {
	app.Post("/signup", auth.Signup)
	app.Post("/login", auth.Login)
	app.Get("/getmydata", base.GetMyData)
	app.Post("/logoutall", auth.LogoutAll)
}
