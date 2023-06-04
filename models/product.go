package models

import (
	"errors"
	"fmt"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title     string `json:"title"`
	Price     uint   `json:"price"`
	Inventory uint   `json:"inventory"`
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
	var products []Product
	db.Find(&products, id)
	return c.JSON(products)
}

func NewProducts(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(User)
	fmt.Println(user.ID)
	db := database.DBConn
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).SendString("Invalid body!")
	}
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
