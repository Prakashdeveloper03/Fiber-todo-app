package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

type Todo struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Complete bool
}

func main() {
	// Set up app config
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Set up database connection
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	// Routes
	app.Get("/", home)
	app.Post("/add", add)
	app.Get("/update/:id", update)
	app.Get("/delete/:id", delete)

	// Start server
	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func home(c *fiber.Ctx) error {
	var todoList []Todo
	db.Find(&todoList)
	return c.Render("base", fiber.Map{
		"todoList": todoList,
	})
}

func add(c *fiber.Ctx) error {
	title := c.FormValue("title")
	newTodo := Todo{Title: title, Complete: false}
	db.Create(&newTodo)
	return c.Redirect("/")
}

func update(c *fiber.Ctx) error {
	var todo Todo
	id := c.Params("id")
	db.First(&todo, id)
	todo.Complete = !todo.Complete
	db.Save(&todo)
	return c.Redirect("/")
}

func delete(c *fiber.Ctx) error {
	var todo Todo
	id := c.Params("id")
	db.Delete(&todo, id)
	return c.Redirect("/")
}
