package models

import (
	"errors"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title     string `json:"title"`
	Price     uint   `json:"price"`
	Inventory uint   `json:"inventory"`
	Image     string `json:"image"`
}

func GetProducts(c *fiber.Ctx) error {
	db := database.DBConn
	var products []Product
	db.Find(&products)
	return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	result := db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}
	return c.JSON(product)
}

func NewProducts(c *fiber.Ctx) error {
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).SendString("Invalid body!")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("Failed to process the form")
	}

	files := form.File["image"]

	// Check the number of files
	if len(files) == 0 {
		return c.Status(400).SendString("An image file is required")
	} else if len(files) > 1 {
		return c.Status(400).SendString("Only one image file is allowed")
	}

	// Get the file from the files slice
	file := files[0]

	if file != nil {
		// Save the uploaded file
		savePath, err := saveUploadedFile(file)
		if err != nil {
			return c.Status(500).SendString("Failed to save image file!" + err.Error())
		}
		product.Image = savePath
	} else {
		return c.Status(500).SendString("Failed to get image file!")
	}

	db := database.DBConn
	db.Create(&product)
	return c.JSON(product)
}

func DeleteProducts(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var product Product
	db.First(&product, id)
	result := db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}
	db.Delete(&product)
	return c.SendString("Deleted")
}
