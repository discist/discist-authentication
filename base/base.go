package base

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	if userID == "" {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utils.NewJError(utils.ErrNoSession))
	}

	if userID != "" {

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

	normailizedusername := utils.NormalizeEmail(updatedata.Username)

	lenght := len(normailizedusername)

	if updatedata.Username != "" && lenght < 18 {

		id, err := controllers.RedisGetKey(updatedata.Uuid)

		if err != nil {

			return ctx.
				Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))
		}

		exist, err := controllers.UGetByKey("username", normailizedusername)
		if err != nil {
			logrus.Error(err)
		}

		requestinguser, err := controllers.UGetByID(id)
		if err != nil {
			logrus.Error(err)
		}

		if exist.Username != "" {

			if exist.Username != requestinguser.Username {
				fmt.Println("EXIST:", exist.Username, "REQ:", requestinguser.Username)
				return ctx.
					Status(http.StatusBadRequest).
					JSON(fiber.Map{"error": "usernametaken"})
			}

		}

		err = controllers.AddNewKey(id, "username", normailizedusername)
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

func UploadProfilePicture(ctx *fiber.Ctx) error {

	var u models.UpdateUser
	err := ctx.BodyParser(&u)

	suuid := u.Uuid
	suuid = strings.Trim(suuid, "\"")

	if err != nil {

		return ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	file, err := ctx.FormFile("attachment")

	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utils.NewJError(err))

	} else {

		id, err := controllers.RedisGetKey(suuid)

		if err != nil {
			return ctx.
				Status(http.StatusUnprocessableEntity).
				JSON(utils.NewJError(err))

		}
		uuidWithHyphen := uuid.New()
		uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -10)
		fmt.Println(uuid)

		filedir := "uploads"
		filename := fmt.Sprintf("%s-%s", suuid, uuid)

		finalstring := fmt.Sprintf("./%s/%s.jpeg", filedir, filename)

		ctx.SaveFile(file, finalstring)
		fmt.Println(file.Size)
		profilephotourl := fmt.Sprintf("http://112.133.192.241:8089/storage/%s-%s.jpeg", suuid, uuid)
		//fmt.Println(profilephotourl)

		err = controllers.AddNewKey(id, "profilephotourl", profilephotourl)
		utils.CheckErorr(err)

		ctx.
			Status(http.StatusAccepted).
			JSON(fiber.Map{"msg": "success", "profilephotourl": profilephotourl})

	}

	// id := "23"
	// add := 23

	// profilepiclink := fmt.Sprintf("%s%s%duserprofilepicture.jpeg", id, "upno", add)

	return err

}

func GetAllUser(ctx *fiber.Ctx) error {

	var alldoc []models.UserAllDataPublic

	alldoc, err := controllers.GetAll()
	if err != nil {
		logrus.Error(err)
	}

	//logrus.Info(alldoc)
	//fmt.Println(alldoc)

	return ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"alluserdata": alldoc})

}
