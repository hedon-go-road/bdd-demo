package routes

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/hedon-go-road/bdd-demo/database"
	"github.com/hedon-go-road/bdd-demo/models"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/books", AddBook)
	app.Get("/books", Book)
	app.Put("/books/:id", Update)
	app.Delete("/books/:id", Delete)
}

func AddBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	database.DB.DB.Create(&book)
	resp := make([]models.Book, 0)
	resp = append(resp, *book)
	return c.Status(http.StatusCreated).JSON(resp)
}

func Book(c *fiber.Ctx) error {
	books := []models.Book{}
	title := c.Query("title")

	if title != "" {
		database.DB.DB.Where("title = ?", title).Find(&books)
	} else {
		database.DB.DB.Find(&books)
	}

	return c.Status(http.StatusOK).JSON(books)
}

func Update(c *fiber.Ctx) error {
	book := new(models.Book)

	id := c.Params("id")
	ui64, _ := strconv.ParseUint(id, 10, 64)
	if err := c.BodyParser(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}
	book.ID = uint(ui64)
	database.DB.DB.Updates(&book)

	return c.Status(http.StatusOK).JSON(book)
}

func Delete(c *fiber.Ctx) error {
	book := []models.Book{}
	title := new(models.Book)
	if err := c.ParamsParser(title); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	database.DB.DB.Where("title = ?", title.Title).Delete(&book)
	return c.Status(http.StatusOK).JSON("deleted")
}
