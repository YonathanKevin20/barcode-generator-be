package handlers

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

var supplierRepo models.SupplierRepository = &models.GormSupplierRepository{}

func SetSupplierRepository(repo models.SupplierRepository) {
	supplierRepo = repo
}

func GetSuppliers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	code := c.Query("code")
	name := c.Query("name")

	offset := (page - 1) * limit

	var filter models.SupplierFilter
	filter.Code = code
	filter.Name = name
	filter.Offset = offset
	filter.Limit = limit

	suppliers, total, err := supplierRepo.FindAllWithFilter(&filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve suppliers")
	}

	return utils.PaginatedResponse(c, suppliers, total, page, limit)
}

func GetSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	supplier, err := supplierRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Supplier not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, supplier)
}

func CreateSupplier(c *fiber.Ctx) error {
	supplier := new(models.Supplier)
	if err := sonic.Unmarshal(c.Body(), &supplier); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	if supplier.Code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Supplier Code is required")
	}
	if supplier.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Supplier Name is required")
	}
	if len(supplier.Code) != 4 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Supplier Code must be exactly 4 characters long")
	}

	if existing, err := supplierRepo.FindByCodeAndName(supplier.Code, supplier.Name); err == nil && existing != nil && existing.ID != 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Supplier with this Code and Name already exists")
	}

	if err := supplierRepo.Create(supplier); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create supplier")
	}
	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "Supplier created successfully"})
}

func UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	supplier := new(models.Supplier)
	if err := sonic.Unmarshal(c.Body(), &supplier); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	existingSupplier, err := supplierRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Supplier not found")
	}

	supplier.ID = existingSupplier.ID // ensure correct ID
	if err := supplierRepo.Update(supplier); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update supplier")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Supplier updated successfully"})
}

func DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	supplier, err := supplierRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Supplier not found")
	}
	if err := supplierRepo.Delete(supplier); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Supplier deleted successfully"})
}
