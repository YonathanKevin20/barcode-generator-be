package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"barcode-generator-be/models"

	"github.com/gofiber/fiber/v2"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) FindAll() ([]models.User, error) {
	var result []models.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockUserRepo) FindAllWithFilter(filter *models.UserFilter) ([]models.User, int64, error) {
	var result []models.User
	for _, u := range m.users {
		if (filter.Username == "" || u.Username == filter.Username) && (filter.Role == "" || string(u.Role) == filter.Role) {
			result = append(result, *u)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockUserRepo) FindByID(id uint) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fiber.ErrNotFound
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, fiber.ErrNotFound
}

func (m *mockUserRepo) Create(user *models.User) error {
	if _, exists := m.users[user.Username]; exists {
		return fiber.ErrConflict
	}
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) Update(user *models.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return fiber.ErrNotFound
	}
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) Delete(user *models.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return fiber.ErrNotFound
	}
	delete(m.users, user.Username)
	return nil
}

func TestCreateUser_MissingUsername(t *testing.T) {
	app := fiber.New()
	app.Post("/users", CreateUser)

	payload := map[string]any{
		"password": "testpass",
		"role":     models.RoleAdmin,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateUser_MissingPassword(t *testing.T) {
	app := fiber.New()
	app.Post("/users", CreateUser)

	payload := map[string]any{
		"username": "testuser",
		"role":     models.RoleAdmin,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	app := fiber.New()
	app.Post("/users", CreateUser)

	payload := map[string]any{
		"username": "testuser",
		"password": "testpass",
		"role":     "invalidrole",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	mockRepo := newMockUserRepo()
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Post("/users", CreateUser)

	payload := map[string]any{
		"username": "duplicateuser",
		"password": "testpass",
		"role":     models.RoleAdmin,
	}
	body, _ := json.Marshal(payload)

	// First request should succeed
	req1 := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp1.StatusCode)
	}

	// Second request with same username should fail
	req2 := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != 409 {
		t.Errorf("expected status 409, got %d", resp2.StatusCode)
	}
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := newMockUserRepo()
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Post("/users", CreateUser)

	payload := map[string]any{
		"username": "newuser",
		"password": "securepass",
		"role":     models.RoleAdmin,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestGetUsers_Empty(t *testing.T) {
	mockRepo := newMockUserRepo()
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Get("/users", GetUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetUsers_WithUsers(t *testing.T) {
	mockRepo := newMockUserRepo()
	mockRepo.Create(&models.User{ID: 1, Username: "user1", Password: "pass", Role: models.RoleAdmin})
	mockRepo.Create(&models.User{ID: 2, Username: "user2", Password: "pass", Role: models.RoleOperator})
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Get("/users", GetUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := newMockUserRepo()
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Get("/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/users/123", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetUser_Found(t *testing.T) {
	mockRepo := newMockUserRepo()
	mockRepo.Create(&models.User{ID: 42, Username: "findme", Password: "pass", Role: models.RoleAdmin})
	SetUserRepository(mockRepo)
	app := fiber.New()
	app.Get("/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/users/42", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
