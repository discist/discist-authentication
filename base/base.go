package base

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetMyData(ctx *fiber.Ctx) error {

	//to discis

	//header := ctx.GetRespHeader("Content-Type")
	authheader := ctx.GetRespHeader("Content-Type")
	fmt.Println(authheader)

	var sessionid models.Session
	err := ctx.BodyParser(&sessionid)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	userID, err := controllers.RedisGetKey(sessionid.Uuid)
	if userID == "" {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))
	}

	if userID != "" {
		//fmt.Println(userID)
		var userInfo models.User
		userInfo, err := controllers.GetByID(userID)
		if err != nil {
			return ctx.
				Status(http.StatusBadGateway).
				JSON(utils.NewJError(utils.ErrNoSession))
		}

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"userdata": userInfo, "headers": authheader})

	}

	return err

}
