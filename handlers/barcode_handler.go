package handlers

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func GetBarcodes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	statusID := c.Query("status_id")
	categoryID := c.Query("category_id")
	categoryName := c.Query("category_name")
	supplierID := c.Query("supplier_id")
	supplierName := c.Query("supplier_name")
	productName := c.Query("product_name")
	barcode := c.Query("barcode")

	offset := (page - 1) * limit

	var filter models.BarcodeFilter
	if statusID != "" {
		if v, err := strconv.Atoi(statusID); err == nil {
			filter.StatusID = uint(v)
		}
	}
	if categoryID != "" {
		if v, err := strconv.Atoi(categoryID); err == nil {
			filter.CategoryID = uint(v)
		}
	}
	if supplierID != "" {
		if v, err := strconv.Atoi(supplierID); err == nil {
			filter.SupplierID = uint(v)
		}
	}
	filter.CategoryName = strings.ToUpper(categoryName)
	filter.SupplierName = strings.ToUpper(supplierName)
	filter.ProductName = strings.ToUpper(productName)
	filter.Barcode = barcode
	filter.Offset = offset
	filter.Limit = limit

	results, total, err := models.BarcodeRepo.FindAllWithFilter(&filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barcodes")
	}

	return utils.PaginatedResponse(c, results, total, page, limit)
}

type ActiveBarcodesQuery struct {
	IDs []uint `query:"id"`
}

func GetActiveBarcodes(c *fiber.Ctx) error {
	var q ActiveBarcodesQuery
	if err := c.QueryParser(&q); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}
	if len(q.IDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No IDs provided")
	}
	barcodes, err := models.BarcodeRepo.FindActiveByIDs(q.IDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barcodes")
	}
	if len(barcodes) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No barcodes found for the provided IDs")
	}
	return utils.JSONResponse(c, fiber.StatusOK, barcodes)
}

func GetBarcode(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	barcode, err := models.BarcodeRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Barcode not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, barcode)
}

func generateBarcode(statusID uint, categoryCode, supplierCode, productCode string) string {
	sequenceID := strconv.Itoa(int(statusID)) + categoryCode + supplierCode + productCode

	// Ensure we have exactly 12 digits for EAN-13 (before checksum)
	if len(sequenceID) < 12 {
		sequenceID = sequenceID + strings.Repeat("0", 12-len(sequenceID))
	} else if len(sequenceID) > 12 {
		sequenceID = sequenceID[:12]
	}

	var oddSum, evenSum int
	for i := range 12 {
		digit, _ := strconv.Atoi(string(sequenceID[i]))
		if i%2 == 0 {
			evenSum += digit
		} else {
			oddSum += digit
		}
	}
	checksum := (10 - ((oddSum*3 + evenSum) % 10)) % 10
	return sequenceID + strconv.Itoa(checksum)
}

func CreateBarcode(c *fiber.Ctx) error {
	barcode := new(models.Barcode)
	if err := sonic.Unmarshal(c.Body(), &barcode); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	_, err := models.StatusRepo.FindByID(barcode.StatusID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Status ID")
	}

	category, err := models.CategoryRepo.FindByID(barcode.CategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Category ID")
	}

	supplier, err := models.SupplierRepo.FindByID(barcode.SupplierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Supplier ID")
	}

	// Check if barcode already exists
	exists, err := models.BarcodeRepo.FindExists(barcode.StatusID, barcode.CategoryID, barcode.SupplierID, strings.ToUpper(barcode.ProductName))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check existing barcode")
	}
	if exists {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Barcode already exists for this product")
	}

	// Generate next product code
	nextProductCode, err := models.BarcodeRepo.GetNextProductCode(barcode.CategoryID, barcode.SupplierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate product code")
	}

	// Generate barcode with the new productCode
	barcode.Barcode = generateBarcode(barcode.StatusID, category.Code, supplier.Code, nextProductCode)

	// Set CreatedByID
	userID := c.Locals("id").(uint)
	barcode.CreatedBy = &userID

	if err := models.BarcodeRepo.Create(barcode); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create barcode")
	}

	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "Barcode created successfully"})
}

func UpdateBarcodeInactive(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	barcode, err := models.BarcodeRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Barcode not found")
	}
	if err := models.BarcodeRepo.UpdateInactive(barcode.ID, !barcode.IsInactive); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update barcode inactive")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Barcode inactive updated successfully"})
}

func DeleteBarcode(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	barcode, err := models.BarcodeRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Barcode not found")
	}
	if err := models.BarcodeRepo.Delete(barcode.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete barcode")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Barcode deleted successfully"})
}
