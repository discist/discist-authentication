package security

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// fmt.Println(hashed)

	return string(hashed), nil
}

func VerifyPassword(hashed, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

}

func SessionLookup(session models.Session, ctx *fiber.Ctx) (models.User, error) {
	var userInfo models.User
	userID, err := controllers.RedisGetKey(session.Uuid)
	fmt.Println(err)
	if userID == "" {
		return userInfo, ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))

	}

	if userID != "" {
		//fmt.Println(userID)

		userInfo, err = controllers.GetByID(userID)
		if err != nil {
			return userInfo, ctx.
				Status(http.StatusBadGateway).
				JSON(utils.NewJError(utils.ErrNoSession))
		}

	}
	return userInfo, ctx.
		Status(http.StatusAccepted).
		JSON(fiber.Map{"userdata": userInfo})
}
