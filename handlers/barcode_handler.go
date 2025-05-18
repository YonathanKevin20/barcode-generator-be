package handlers

import (
	"barcode-generator-be/config"
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

var barcodeRepo models.BarcodeRepository = &models.GormBarcodeRepository{}

func SetBarcodeRepository(repo models.BarcodeRepository) {
	barcodeRepo = repo
}

func GetBarcodes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	statusID := c.Query("status_id")
	barcode := c.Query("barcode")

	offset := (page - 1) * limit

	var filter models.BarcodeFilter
	if statusID != "" {
		if v, err := strconv.Atoi(statusID); err == nil {
			filter.StatusID = uint(v)
		}
	}
	filter.Barcode = barcode
	filter.Offset = offset
	filter.Limit = limit

	results, total, err := barcodeRepo.FindAllWithFilter(&filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barcodes")
	}

	return utils.PaginatedResponse(c, results, total, page, limit)
}

func GetBarcode(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	barcode, err := barcodeRepo.FindByID(uint(idUint))
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

	// Check if Status exists
	var status models.Status
	if err := config.DB.First(&status, "id = ?", barcode.StatusID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Status ID")
	}

	// Check if Category exists
	var category models.Category
	if err := config.DB.First(&category, "id = ?", barcode.CategoryID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Category ID")
	}

	// Check if Supplier exists
	var supplier models.Supplier
	if err := config.DB.First(&supplier, "id = ?", barcode.SupplierID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Supplier ID")
	}

	// Check if Barcode already exists
	var existingBarcode models.Barcode
	if err := config.DB.First(&existingBarcode, "status_id = ? AND category_id = ? AND supplier_id = ? AND UPPER(product_name) = ?", barcode.StatusID, barcode.CategoryID, barcode.SupplierID, strings.ToUpper(barcode.ProductName)).Error; err == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Barcode already exists")
	}

	// Generate next product code
	nextProductCode, err := barcodeRepo.GetNextProductCode(barcode.CategoryID, barcode.SupplierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate product code")
	}

	// Generate barcode with the new productCode
	barcode.Barcode = generateBarcode(barcode.StatusID, category.Code, supplier.Code, nextProductCode)

	if err := config.DB.Create(&barcode).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create barcode")
	}

	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "Barcode created successfully"})
}

func DeleteBarcode(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	barcode, err := barcodeRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Barcode not found")
	}
	if err := barcodeRepo.Delete(barcode.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete barcode")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Barcode deleted successfully"})
}
