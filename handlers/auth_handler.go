package handlers

import (
	"barcode-generator-be/config"
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	input := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := sonic.Unmarshal(c.Body(), &input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	if input.Username == "" || input.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username and password are required")
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not hash password")
	}

	user := &models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "User created successfully"})
}

func Login(c *fiber.Ctx) error {
	input := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	if err := sonic.Unmarshal(c.Body(), &input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	if input.Username == "" || input.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username and password are required")
	}

	user := new(models.User)
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate token")
	}

	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"token": token})
}

func GetMe(c *fiber.Ctx) error {
	id := c.Locals("id").(uint)

	user := new(models.User)
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.JSONResponse(c, fiber.StatusOK, user)
}

func Logout(c *fiber.Ctx) error {
	jti := c.Locals("jti").(string)
	claims, err := utils.GetTokenClaims(c.Get("Authorization")[7:]) // Remove "Bearer " prefix
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not process token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid token expiration")
	}

	expirationTime := time.Unix(int64(exp), 0)
	err = utils.TokenBlacklist.Add(jti, expirationTime)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not blacklist token")
	}

	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "Successfully logged out"})
}
