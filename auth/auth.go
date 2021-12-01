package auth

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/security"
	"authentication/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/asaskevich/govalidator.v9"
)

func Signup(ctx *fiber.Ctx) error {

	var newuser models.User

	err := ctx.BodyParser(&newuser)
	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utils.NewJError(err))
	}

	newuser.Email = utils.NormalizeEmail(newuser.Email)

	if !govalidator.IsEmail(newuser.Email) {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrInvalidEmail))

	}

	exist, err := controllers.GetByEmail(newuser.Email)
	utils.CheckErorr(err)

	if exist.Email != "" {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrEmailAlreadyExists))
	}

	if exist.Email == "" {

		if strings.TrimSpace(newuser.Password) == "" {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(utils.ErrEmptyPassword))

		}

		if len(newuser.Password) < 5 {
			fmt.Println("short password")
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(utils.ErrShortPassword))
		}

		newuser.Password, err = security.EncryptPassword(newuser.Password)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))

		}

		newuser.CreatedAt = time.Now()
		newuser.UpdatedAt = newuser.CreatedAt
		newuser.ID = primitive.NewObjectID()
		newuser.Sessions = []models.Session{}

		err := controllers.Save(&newuser)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"user created": "gg"})

	}

	return err

}

func Login(ctx *fiber.Ctx) error {
	var loginUser models.User

	err := ctx.BodyParser(&loginUser)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utils.NewJError(err))
	}

	loginUser.Email = utils.NormalizeEmail(loginUser.Email)
	if !govalidator.IsEmail(loginUser.Email) {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrInvalidEmail))

	}

	userinfo, err := controllers.GetByEmail(loginUser.Email)
	utils.CheckErorr(err)
	fmt.Println(userinfo.Email)
	if err != nil {
		log.Println("login failed")
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrIncorrectEmail))

	}

	err = security.VerifyPassword(userinfo.Password, loginUser.Password)
	if err == nil {
		fmt.Printf("creating session from user %s \n", userinfo.ID)

		stringObjectID := userinfo.ID.Hex()
		uuidh := uuid.New()
		uuid := strings.Replace(uuidh.String(), "-", "", -1)

		controllers.RedisAddKey(uuid, stringObjectID)
		fmt.Println("Added to redis")
		var sessiondata models.Session
		sessiondata.Uuid = uuid
		sessiondata.Device = "macbook pro"
		sessiondata.Location = "cloudflake hong-kong"
		var sessiondatas models.Session
		sessiondatas.Uuid = uuid
		sessiondatas.Device = "iphone 11 pro"
		sessiondatas.Location = "act fiber chennei"
		sessionarray := []models.Session{sessiondata, sessiondatas}

		var updateduserinfo models.User

		updateduserinfo.Sessions = sessionarray

		controllers.UpdateSessions("_id", stringObjectID, updateduserinfo)

		return ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"uuid": uuid})

	}
	if err != nil {

		log.Println(err)

		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrIncorrectPassword))

	}

	return err

}
