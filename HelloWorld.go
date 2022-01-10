package main

import "github.com/gofiber/fiber/v2"

func main() {
  app := fiber.New()

//   app.Get("/", func(c *fiber.Ctx) error {
//     return c.SendString("Hello, World!")
//   })
  app.Get("/api/todo", func(c *fiber.Ctx) error {
     return c.SendString("Hello, World!")
   })
  app.Post("/api/item", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })
  app.Delete("/api/item/:id", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })
  app.Patch("/api/item/:id", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })

  app.Listen(":3000")
  //asdß
}