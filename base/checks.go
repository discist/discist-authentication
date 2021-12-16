package base

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckUserNameAvail(ctx *fiber.Ctx) error {

	var u models.UsernameCheck

	err := ctx.BodyParser(&u)

	username := utils.NormalizeEmail(u.Username)

	lenght := len(username)

	fmt.Println(username)
	if err != nil {

		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrEmailAlreadyExists))

	}

	_, err = controllers.GetByKey("username", username)
	if err != nil && lenght < 18 {

		if err == mongo.ErrNoDocuments {
			return ctx.Status(http.StatusAccepted).
				JSON(fiber.Map{"avail": true, "username": username})
		}

		return ctx.Status(http.StatusAccepted).
			JSON(fiber.Map{"avail": false, "username": username})

	}

	return ctx.
		Status(http.StatusAccepted).
		JSON(fiber.Map{"avail": false, "username": username})

}
