package todos

import (
    // import modules
    "strconv"

    // import packages
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "github.com/jindongh/receipt/database"
)

type Todo struct {
    gorm.Model
    Id int `gorm:"primaryKey"`
    Name string `json:"name"`
    Completed bool `json:"completed"`
}

func GetAll(c *fiber.Ctx) error {
    db := database.DBConn
    var todoss []Todo
    db.Find(&todoss)
    // If the database read is successful
    return c.Status(fiber.StatusOK).JSON(todoss)
}

func GetOne(ctx *fiber.Ctx) error {
    paramsId := ctx.Params("id")
    id, err := strconv.Atoi(paramsId)
    if err != nil {
        ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse id",
        })
        return err
    }

    db := database.DBConn

    var todo Todo
    db.Find(&todo, id)

    // If the database read is successful
    if int(todo.Id) == id{
        return ctx.Status(fiber.StatusOK).JSON(todo)
    }

    // If the database fails to read the id parameter
    return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error": "todo not found",
    })
}

func AddTodo(ctx *fiber.Ctx) error {
    db := database.DBConn
    type request struct {
        Name string `json:"name"`
    }
    // Parse POST data
    var body request
    err := ctx.BodyParser(&body)
    if err != nil {
        ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse json",
        })
        return err
    }
    // Get the json struct that is required to send
    id := uuid.New()
    todo := Todo{
        Id: int(id.ID()),
        Name: body.Name,
        Completed: false,
        }
    // Insert to DB
    db.Create(&todo)

    return ctx.Status(fiber.StatusOK).JSON(todo)
}

func DeleteTodo(ctx *fiber.Ctx) error {
    db := database.DBConn
    paramsId := ctx.Params("id")
    id, err := strconv.Atoi(paramsId)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse id",
        })
    }

    var todo Todo
    db.First(&todo, id)

    if int(todo.Id) != id {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "todo not found",
        })
    }

    db.Delete(&todo)

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "todo deleted successfully",
    })
}

func UpdateTodo(ctx *fiber.Ctx) error {
    db := database.DBConn

    type request struct {
        Name *string `json:"name"`
        Completed *bool `json:"completed"`
    }

    paramsId := ctx.Params("id")
    id, err := strconv.Atoi(paramsId)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse id",
        })
    }

    var body request

    err = ctx.BodyParser(&body)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error" : "Cannot parse body",
        })
    }

    var todo Todo
    // Check if todo exist, if exist assign it value to todo
    db.First(&todo, id)

    // handling 404 error
    if int(todo.Id) != id {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "todo not found",
        })
    }

    if body.Name != nil {
        todo.Name = *body.Name
    }

    if body.Completed != nil {
        todo.Completed = *body.Completed
    }

    db.Save(&todo)

    return ctx.Status(fiber.StatusOK).JSON(todo)
}

