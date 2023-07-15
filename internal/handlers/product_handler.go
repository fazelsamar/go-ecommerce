package handlers

import (
	"errors"
	"strconv"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler interface {
	GetProductList(*fiber.Ctx) error
	GetProduct(*fiber.Ctx) error
	NewProduct(*fiber.Ctx) error
	DeleteProduct(*fiber.Ctx) error
}

type productHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) ProductHandler {
	return productHandler{
		productService: productService,
	}
}

func (ph productHandler) GetProductList(c *fiber.Ctx) error {
	products, err := ph.productService.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Cant get product list: " + err.Error()})
	}
	return c.Status(200).JSON(products)
}

func (ph productHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return c.Status(500).JSON(fiber.Map{"massage": "Invalid id"})
	}

	product, err := ph.productService.GetProductById(uint(id))
	if err != nil || product.ID == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"massage": "Not found"})
		} else {
			return c.Status(500).JSON(fiber.Map{"massage": "Something went wrong"})
		}
	}
	return c.Status(200).JSON(product)
}

func (ph productHandler) NewProduct(c *fiber.Ctx) error {
	product := new(entity.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Invalid body: " + err.Error()})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Failed to process the form"})
	}

	files := form.File["image"]

	// Check the number of files
	if len(files) == 0 {
		return c.Status(500).JSON(fiber.Map{"massage": "An image file is required"})
	} else if len(files) > 1 {
		return c.Status(500).JSON(fiber.Map{"massage": "Only one image file is allowed"})
	}

	// Get the file from the files slice
	file := files[0]

	if file != nil {
		// Save the uploaded file
		savePath, err := saveUploadedFile(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"massage": "Failed to save image file!" + err.Error()})
		}
		product.Image = savePath
	} else {
		return c.Status(500).JSON(fiber.Map{"massage": "Failed to get image file!"})
	}

	product, err = ph.productService.CreateProduct(product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Can't create product: " + err.Error()})
	}

	return c.JSON(product)
}

func (ph productHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return c.Status(500).JSON(fiber.Map{"massage": "Invalid id"})
	}

	product, err := ph.productService.GetProductById(uint(id))
	if err != nil || product.ID == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"massage": "Not found"})
		} else {
			return c.Status(500).JSON(fiber.Map{"massage": "Something went wrong"})
		}
	}
	err = ph.productService.DeleteProduct(&product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Can't delete product"})
	}

	return c.Status(200).JSON(fiber.Map{"massage": "Successfully deleted"})
}
