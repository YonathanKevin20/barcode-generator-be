package handlers

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	code := c.Query("code")
	name := c.Query("name")

	offset := (page - 1) * limit

	var filter models.CategoryFilter
	filter.Code = code
	filter.Name = name
	filter.Offset = offset
	filter.Limit = limit

	categories, total, err := models.CategoryRepo.FindAllWithFilter(&filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve categories")
	}

	return utils.PaginatedResponse(c, categories, total, page, limit)
}

func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	category, err := models.CategoryRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, category)
}

func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := sonic.Unmarshal(c.Body(), &category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	if category.Code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category Code is required")
	}
	if category.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category Name is required")
	}
	if len(category.Code) != 4 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category Code must be exactly 4 characters long")
	}

	if existing, err := models.CategoryRepo.FindByCode(category.Code); err == nil && existing != nil && existing.ID != 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Category with this Code already exists")
	}

	if existing, err := models.CategoryRepo.FindByCodeAndName(category.Code, category.Name); err == nil && existing != nil && existing.ID != 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Category with this Code and Name already exists")
	}

	if err := models.CategoryRepo.Create(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create category")
	}
	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "Category created successfully"})
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	category := new(models.Category)
	if err := sonic.Unmarshal(c.Body(), &category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	existingCategory, err := models.CategoryRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}

	category.ID = existingCategory.ID // ensure correct ID
	if err := models.CategoryRepo.Update(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update category")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Category updated successfully"})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	category, err := models.CategoryRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}
	if err := models.CategoryRepo.Delete(category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete category")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Category deleted successfully"})
}
