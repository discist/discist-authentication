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

	var sessionid models.Session

	err := ctx.BodyParser(&sessionid)
	if sessionid.Uuid == "" {

		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrEmailAlreadyExists))

	}
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	userID, err := controllers.RedisGetKey(sessionid.Uuid)
	fmt.Println("sesssionid from frontend ", sessionid)
	if userID == "" {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))
	}

	if userID != "" {
		fmt.Println(userID, "this is obj id")
		var userInfo models.UserAllData
		userInfo, err := controllers.GetFullDoc(userID)
		if err != nil {
			return ctx.
				Status(http.StatusBadGateway).
				JSON(utils.NewJError(utils.ErrNoSession))
		}

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"userdata": userInfo})

	}

	return err

}

func UpdateData(ctx *fiber.Ctx) error {

	var updatedata models.UpdateUser

	err := ctx.BodyParser(&updatedata)

	fmt.Println(updatedata)

	if updatedata.Username != "" {

		id, err := controllers.RedisGetKey(updatedata.Uuid)
		if err != nil {

			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}

		err = controllers.AddNewKey(id, "username", updatedata.Username)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}
		err = controllers.AddNewKey(id, "story", updatedata.Story)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}
		err = controllers.AddNewKey(id, "subject", updatedata.Subject)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}

		err = controllers.AddNewKey(id, "state", updatedata.State)

		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"response": "data updated"})
	}

	return err

}
