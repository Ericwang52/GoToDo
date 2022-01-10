package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
  )
  func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
type User struct {
	ID           uint
	Name         string
	Password 	 string
  }
func main() {
	dsn := "host=localhost user=ericwang dbname=ericwang port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  db.AutoMigrate(&User{})

	user := User{Name: "test", Password: "test"}

	db.Create(&user) // pass pointer of data to Create

  app := fiber.New()

//   app.Get("/", func(c *fiber.Ctx) error {
//     return c.SendString("Hello, World!")
//   })
  app.Get("/api/items", func(c *fiber.Ctx) error {
    result := map[string]interface{}{}
      db.Table("users").Take(&result)
     return c.SendString(result["name"].(string))
   })
  app.Post("/api/items", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })
  app.Delete("/api/items/:id", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })
  app.Patch("/api/items/:id", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
  })

  app.Listen(":3000")
  //asdß
}