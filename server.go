package main

import (
  "log"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/cors"
  "gofiber-starter/routes"
)

func main() {
  app := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    // AllowHeaders:  "Origin, Content-Type, Accept",
  }))

  routes.Load(app)

  app.Get("/*", func(ctx *fiber.Ctx) error {
    return ctx.SendString("Not found.")
  })

  log.Fatal(app.Listen(":3000"))
}
