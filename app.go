package main

import (
	"embed"
	"html/template"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2" 
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jindongh/receipt/todos"
	"github.com/jindongh/receipt/database"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func initDatabase() {
    var err error
    dbUrl := os.Getenv("DB_URL")
    database.DBConn, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    log.Println("Database successfully connected")

    database.DBConn.AutoMigrate(&todos.Todo{})
    log.Println("Database Migrated")
}

func setupV1(app *fiber.App) {
    // Group is used for Routes with a common prefix to define a new sub-router with optional middleware.
    v1 := app.Group("/v1")
    //Each route will have /v1 prefix
    setupTodosRoutes(v1)
}

func setupTodosRoutes(grp fiber.Router) {
    // Group is used for Routes with common prefix => Each route will have /todos prefix
    todosRoutes := grp.Group("/todos")
    // Route for Get all todos -> navigate to => http://127.0.0.1:3000/v1/todos/
    todosRoutes.Get("/", todos.GetAll)
    // Route for Get a todo -> navigate to => http://127.0.0.1:3000/v1/todos/<todo's id>
    todosRoutes.Get("/:id", todos.GetOne)
    // Route for Add a todo -> navigate to => http://127.0.0.1:3000/v1/todos/
    todosRoutes.Post("/", todos.AddTodo)
    // Route for Delete a todo -> navigate to => http://127.0.0.1:3000/v1/todos/<todo's id>
    todosRoutes.Delete("/:id", todos.DeleteTodo)
    // Route for Update a todo -> navigate to => http://127.0.0.1:3000/v1/todos/<todo's id>
    todosRoutes.Patch("/:id", todos.UpdateTodo)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	app := fiber.New()

	initDatabase()
	setupV1(app)

	app.Get("/", func(c *fiber.Ctx) error {
		// send text
		return c.SendString("Hello, World!")
	})

	app.Use(logger.New(logger.Config{ // add Logger middleware with config
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	err := app.Listen(":" + port) 
	if err != nil {
		panic(err)
	}
}
