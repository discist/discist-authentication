package posts

import (
	"authentication/controllers"
	"authentication/models"
	"authentication/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddPost(ctx *fiber.Ctx) error {

	var post models.Post
	//var comments []models.Comment

	// var comet models.Comment
	// comet.Comment = " hey this is cool"
	// comet.Date = time.Now()
	// comet.Interactions = 2
	// comet.Username = "first comment"
	// var cometx models.Comment
	// cometx.Comment = " hey this is cool"
	// cometx.Date = time.Now()
	// cometx.Interactions = 2
	// cometx.Username = "first comment"

	// comments = append(comments, comet)
	// comments = append(comments, cometx)

	// post.Comments = comments

	err := ctx.BodyParser(&post)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}
	var userdata models.UpdateUser
	userid, err := controllers.RedisGetKey(post.Uuid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}
	userdata, err = controllers.UGetByID(userid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}
	post.Date = time.Now()
	post.Uuid = userid
	post.Username = userdata.Username

	poid, err := controllers.PostSave(&post)

	postID := poid.InsertedID.(primitive.ObjectID).Hex()

	//fmt.Println(postID, "post id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	//fmt.Println(userid, "user id")

	//fmt.Println(userdata, " all data of user")

	Existingposts := userdata.Posts
	Existingposts = append(Existingposts, postID)

	controllers.Updateposts("_id", userid, Existingposts)

	ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"msg": "postsuccess", "post": post})

	return err

}

func Getallposts(ctx *fiber.Ctx) error {

	var req models.Request

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).
			JSON(utils.NewJError(err))
	}

	oid, err := controllers.RedisGetKey(req.Uuid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))
	}

	var uuser models.UpdateUser

	uuser, err = controllers.UGetByID(oid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	fmt.Println("allpost requested by: ", uuser.Username)

	var allpost []models.Post

	allpost, err = controllers.GetAllPosts()

	//fmt.Println(allpost)

	ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"allposts": allpost})

	return err

}

func GetMyPosts(ctx *fiber.Ctx) error {

	var u models.Request

	err := ctx.BodyParser(&u)

	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}
	//fmt.Println(u, "tjhis SAWNHJDSAJDH")

	oid, err := controllers.RedisGetKey(u.Uuid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}
	var userdocs models.UpdateUser

	userdocs, err = controllers.UGetByID(oid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	fmt.Println(userdocs.Posts)

	var post models.Post
	var postarray []models.Post

	for _, postID := range userdocs.Posts {
		fmt.Println(postID, " post id array for loop")

		post, err = controllers.GetPostByID(postID)
		if err != nil {
			ctx.Status(http.StatusBadRequest).
				JSON(utils.NewJError(err))

		}

		postarray = append(postarray, post)

	}

	return ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"posts": postarray})

}

func DeletePost(ctx *fiber.Ctx) error {

	var u models.DeletePost

	err := ctx.BodyParser(&u)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	var userdocs models.UpdateUser
	oid, err := controllers.RedisGetKey(u.Uuid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	userdocs, err = controllers.UGetByID(oid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).
			JSON(utils.NewJError(err))

	}

	var newarray []string

	//fmt.Println(userdocs.Posts, " before delete ")
	for _, postID := range userdocs.Posts {

		if postID != u.PostID {
			newarray = append(newarray, postID)
			//fmt.Println("not deleted")

		} else {
			fmt.Println("deleted: ", postID)
			controllers.PostDelete(u.PostID)
		}

	}
	//fmt.Println(newarray, " array after delete ")

	err = controllers.AddNewArray(oid, "postID", newarray)
	if err != nil {
		log.Fatal(err)
	}

	return ctx.Status(http.StatusAccepted).
		JSON(fiber.Map{"msg": "delete successfull", "postarr": newarray})

}
