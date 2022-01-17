package main

import (
  "github.com/gofiber/fiber/v2/middleware/csrf"
  "fmt"
  "github.com/gofiber/fiber/v2/middleware/session"
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
  type UserIn struct {
    UserID uint8
    //jwt here
  
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
  session := session.New()
  // engine := html.New("./views", ".html")

  app := fiber.New()
  app.Use(csrf.New())
app.Get("/", func(c *fiber.Ctx) error {

    // Render index template
    userIn := new(UserIn)
    c.BodyParser(userIn)  
    var user User
    db.Where("id = ?", userIn.UserID).Preload("ToDos").Find(&user)

    return c.Render("index.html", fiber.Map{
    })
})

  app.Get("/api/items", func(c *fiber.Ctx) error {
    // userIn := new(UserIn)
    // c.BodyParser(userIn)  
    var user User
    db.Where("id = ?", 1).Preload("ToDos").Find(&user)
    return c.JSON(user)
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
  
    u := &User{}
    user := new(User)

    err := c.BodyParser(user)


    db.Where("name = ?", user.Name).Where("password = ?", user.Password).Preload("ToDos").Find(&u)
  
    // if errors.Is(err, gorm.ErrRecordNotFound) {
    //   return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
    // }
  
    // if err := password.Verify(u.Password, b.Password); err != nil {
    //   return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
    // }
  
    // t := jwt.Generate(&jwt.TokenPayload{
    //   ID: u.ID,
    // })
  
    sess, err := session.Get(c)
    if err != nil {
      panic(err)
    }
  
    // Set key/value
    sess.Set("user", u.ID)
    // sess.Set("token", t)
  
    // save session
    if err := sess.Save(); err != nil {
      panic(err)
    }
  
    // Expire csrf cookie
    // ctx.Cookie(&fiber.Cookie{
    // 	Name:    "csrf_",
    // 	Expires: time.Now().Add(-1 * time.Minute),
    // })
  
    return c.JSON(u)

})
  app.Post("/api/items", func(c *fiber.Ctx) error {
     userIn := new(UserIn)
     c.BodyParser(userIn)  
     toDo := new(ToDo)
     c.BodyParser(toDo)  
    var user User
    db.Where("id = ?", userIn.UserID).Preload("ToDos").Find(&user)
    db.Model(&user).Association("ToDos").Append(toDo)

    return c.JSON(fiber.Map{"status": "success", "message": "Created ToDo", "data": user})
  })
  app.Delete("/api/items/:id", func(c *fiber.Ctx) error {
    userIn := new(UserIn)
    c.BodyParser(userIn)  
    var todo ToDo
    db.Where("id = ?", c.Params("id")).Where("user_id = ?", userIn.UserID).Find(&todo)
    // var user User
    // db.Where("id = ?",todo.UserID).Find(&user)
    db.Delete(&todo)
    // db.Model(&user).Association("ToDos").Delete(todo)
    var user User
    db.Where("id = ?", 1).Preload("ToDos").Find(&user)
    return c.JSON(fiber.Map{"status": "success", "message": "Deleted ToDo", "data": user})
  })
  app.Patch("/api/items/:id", func(c *fiber.Ctx) error {
    todoIn := new(ToDo)
    userIn := new(UserIn)
    c.BodyParser(userIn)  
    // Store the body in the note and return error if encountered
    c.BodyParser(todoIn)

    var todo ToDo
    fmt.Print(userIn.UserID);
    db.Where("id = ? AND user_id = ?", c.Params("id"),userIn.UserID).Find(&todo)
    todo.Content = todoIn.Content
    todo.Done = todoIn.Done
    db.Save(&todo)
    var user User
    db.Where("id = ?", 1).Preload("ToDos").Find(&user)
	  return  c.JSON(fiber.Map{"status": "success", "message": "Patched ToDo", "data": user})
  })

  app.Listen(":3000")
  //asdß
}