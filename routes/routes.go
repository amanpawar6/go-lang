package routes

import (
	"example.com/controller"
	"github.com/gofiber/fiber/v2"
)

func handleRequests(app *fiber.App) {
	app.Get("/", controller.HomePage)
	app.Post("/test", controller.TestApi)
	app.Get("/users", controller.UsersDetails)
	app.Get("/user", controller.SingleUserDetails)
	app.Post("/createuser", controller.Createuser)
	app.Put("/updateuser", controller.UpdateUser)
	app.Delete("/deleteuser", controller.DeleteUser)
	app.Post("/login", controller.Userlogin)
	app.Put("/updatepassword", controller.UpdatePassword)
	app.Put("/forgotpassword", controller.ForgetPassword)
}

func Routes(app *fiber.App) {
	handleRequests(app)
}
