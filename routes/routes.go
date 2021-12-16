package routes

import (
	"authentication/auth"
	"authentication/base"

	"github.com/gofiber/fiber/v2"
)

func Install(app *fiber.App) {
	app.Post("/signup", auth.Signup)
	app.Post("/login", auth.Login)

	app.Post("/logoutall", auth.LogoutAll)
	app.Post("/logout", auth.Logoutsession)

	app.Post("/checkusernameavailablity", base.CheckUserNameAvail)

	app.Post("/getmydata", base.GetMyData)
	app.Post("/getallusers", base.GetAllUser)

	app.Post("/updatedata", base.UpdateData)
	app.Post("/uploadprofilepicture", base.UploadProfilePicture)

}
