package handlers

import (
	"barcode-generator-be/models"
	"barcode-generator-be/utils"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var userRepo models.UserRepository = &models.GormUserRepository{}

func SetUserRepository(repo models.UserRepository) {
	userRepo = repo
}

func GetUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	usename := c.Query("username")
	role := c.Query("role")

	offset := (page - 1) * limit

	var filter models.UserFilter
	filter.Username = usename
	filter.Role = role
	filter.Offset = offset
	filter.Limit = limit

	users, total, err := userRepo.FindAllWithFilter(&filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	return utils.PaginatedResponse(c, users, total, page, limit)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	user, err := userRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, user)
}

func isValidRole(role models.Role) bool {
	return role == models.RoleAdmin || role == models.RoleOperator
}

func CreateUser(c *fiber.Ctx) error {
	input := new(struct {
		Username string      `json:"username"`
		Password string      `json:"password"`
		Role     models.Role `json:"role"`
	})
	if err := sonic.Unmarshal(c.Body(), &input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	if input.Username == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Username is required")
	}

	if input.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Password is required")
	}

	if !isValidRole(input.Role) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role")
	}

	// Check if user already exists
	existingUser, err := userRepo.FindByUsername(input.Username)
	if err == nil && existingUser != nil {
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
		Role:     input.Role,
	}

	if err := userRepo.Create(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.JSONResponse(c, fiber.StatusCreated, fiber.Map{"message": "User created successfully"})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	input := new(struct {
		Username    string      `json:"username"`
		NewPassword string      `json:"new_password"`
		Role        models.Role `json:"role"`
	})
	if err := sonic.Unmarshal(c.Body(), &input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	existingUser, err := userRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	if input.Username != "" {
		// Check if the new username is already taken by a different user
		potentialUser, err := userRepo.FindByUsername(input.Username)
		if err == nil && potentialUser != nil && potentialUser.ID != existingUser.ID {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Username already exists")
		}
		existingUser.Username = input.Username
	}

	if input.Role != "" {
		if !isValidRole(input.Role) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role")
		}
		existingUser.Role = input.Role
	}

	// If a new password is provided, hash it
	if input.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Could not hash password")
		}
		existingUser.Password = string(hashedPassword)
	}

	// Update the user
	if err := userRepo.Update(existingUser); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user")
	}

	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "User updated successfully"})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint, _ := strconv.ParseUint(id, 10, 32)
	user, err := userRepo.FindByID(uint(idUint))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	if err := userRepo.Delete(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user")
	}
	return utils.JSONResponse(c, fiber.StatusOK, fiber.Map{"message": "User deleted successfully"})
}
