package main

import (
	"authentication/controllers"
	"authentication/routes"
	"fmt"
	"log"

	//fiber
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	fmt.Println("--- AUTHENTICATION SESSIONS ENDPOINT ---")

	defer controllers.Close()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Static("/storage", "./uploads")
	app.Static("/icons", "./icons")
	app.Use(cors.New(cors.Config{
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE",
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "true",
	}))

	routes.Install(app)

	log.Fatal(app.Listen(":3001"))

}
