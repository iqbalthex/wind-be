package routes

import (
  "github.com/gofiber/fiber/v2"
  "gofiber-starter/controller"
)

var datasetController = controller.NewDatasetController()

func Load(app *fiber.App) {
  api := app.Group("/api/v1")

  datasetController.BindRoute(api.Group("/datasets"))
}
