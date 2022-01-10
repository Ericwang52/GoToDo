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
	ID uint8
	Name         string
	Password 	 string
  ToDos []ToDo 

  }
  type ToDo struct {
    ID uint8
    UserID uint8
    Content string
    Done bool
  }
func main() {
	dsn := "host=localhost user=ericwang dbname=ericwang port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  db.AutoMigrate(&User{}, &ToDo{})

  app := fiber.New()


  app.Get("/api/items", func(c *fiber.Ctx) error {
    var users []User
    db.Where("id = ?", 1).Preload("ToDos").Find(&users)
    return c.JSON(users)
    })
  app.Post("/api/users/register", func(c *fiber.Ctx) error {
    user := new(User)

    // Store the body in the note and return error if encountered
    err := c.BodyParser(user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
    }
    user.ToDos=[]ToDo{{Content:"asdd",Done:false}}
    err = db.Create(&user).Error
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create note", "data": err})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "Created User", "data": user})

    // db.Create(&user) // pass pointer of data to Create
    // return c.JSON(fiber.Map{"status": "success", "message": "user created", "data": nil})
  })
  app.Post("/login", func(c *fiber.Ctx) error {


    return c.Redirect("/")
})
  app.Post("/api/items", func(c *fiber.Ctx) error {
    // user := new(User)
    // c.BodyParser(user)  
    var user User
    db.Where("id = ?", 1).Preload("ToDos").Find(&user)
    db.Model(&user).Association("ToDos").Append(&ToDo{Content:"asddd",Done:false})

    return c.JSON(fiber.Map{"status": "success", "message": "Created ToDo", "data": user})
  })
  app.Delete("/api/items/:id", func(c *fiber.Ctx) error {
    var todo ToDo
    db.Where("id = ?", c.Params("id")).Find(&todo)
    // var user User
    // db.Where("id = ?",todo.UserID).Find(&user)
    db.Delete(&todo)
    // db.Model(&user).Association("ToDos").Delete(todo)

    return c.JSON(fiber.Map{"status": "success", "message": "Deleted ToDo", "data": todo})
  })
  app.Patch("/api/items/:id", func(c *fiber.Ctx) error {
    todoIn := new(ToDo)

    // Store the body in the note and return error if encountered
    c.BodyParser(todoIn)

    var todo ToDo
    db.Where("id = ?", c.Params("id")).Find(&todo)
    todo.Content = todoIn.Content
    todo.Done = todoIn.Done
    db.Save(&todo)
	  return c.SendString("Hello, World!")
  })

  app.Listen(":3000")
  //asdß
}