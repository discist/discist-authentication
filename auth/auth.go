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
		fmt.Println("creating session object to add to db")
		var Newsessiondata models.Session
		Newsessiondata.Uuid = uuid
		UserAgent := ctx.GetRespHeader("User-Agent")
		UserIp := ctx.IP()
		deeets := fmt.Sprintf(`IP:%s Device: %s`, UserIp, UserAgent)
		fmt.Println(deeets, "these are the dear deets")
		Newsessiondata.Device = deeets
		Newsessiondata.Location = "cloudflake hong-kong"

		var res models.User
		res, err = controllers.GetByID(userinfo.ID.Hex())
		if err != nil {
			return ctx.
				Status(http.StatusInternalServerError).
				JSON(utils.NewJError(err))
		}
		ExistingSesionsArray := res.Sessions
		ExistingSesionsArray = append(ExistingSesionsArray, Newsessiondata)
		controllers.UpdateSessions("_id", stringObjectID, ExistingSesionsArray)

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

func LogoutAll(ctx *fiber.Ctx) error {

	fmt.Println("deleting all sessions")
	var sessionID models.Session

	err := ctx.BodyParser(&sessionID)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}
	if sessionID.Uuid == "" {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(utils.ErrNoSession))

	}

	userID, err := controllers.RedisGetKey(sessionID.Uuid)
	if userID == "" {
		fmt.Println(err)
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))
	}

	var userInfo models.User

	if userID != "" {
		userInfo, err = controllers.GetByID(userID)
		if err != nil {
			return ctx.
				Status(http.StatusBadGateway).
				JSON(utils.NewJError(utils.ErrNoSession))
		}

	}

	fmt.Println(userInfo)

	empty := []models.Session{}

	err = controllers.UpdateSessions("_id", userInfo.ID.Hex(), empty)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	return ctx.
		Status(http.StatusAccepted).
		JSON(fiber.Map{"deleted-all": userInfo.Sessions})

}
