package base

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetUserData(ctx *fiber.Ctx) error {

	var user models.GetUser

	err := ctx.BodyParser(&user)
	if user.Uuid == "" {

		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrEmailAlreadyExists))

	}
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	userID, err := controllers.RedisGetKey(user.Uuid)

	if userID == "" {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))
	}

	if userID != "" {

		var userInfo models.UserAllDataPublic
		userInfo, err := controllers.AllDataGetByKey("username", user.UserName)
		if err != nil {
			return ctx.
				Status(http.StatusBadGateway).
				JSON(utils.NewJError(utils.ErrNoSession))
		}

		fmt.Println(userInfo)

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"userdata": userInfo})

	}

	return err

}
