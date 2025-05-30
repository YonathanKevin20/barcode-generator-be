package handlers

import (
	"barcode-generator-be/models"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockStatusRepo struct {
	statuses map[string]*models.Status
}

func newMockStatusRepo() *mockStatusRepo {
	return &mockStatusRepo{statuses: make(map[string]*models.Status)}
}

func (m *mockStatusRepo) FindAll() ([]models.Status, error) {
	var result []models.Status
	for _, s := range m.statuses {
		result = append(result, *s)
	}
	return result, nil
}

func (m *mockStatusRepo) FindByID(id uint) (*models.Status, error) {
	for _, u := range m.statuses {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fiber.ErrNotFound
}

func TestGetStatuses_Success(t *testing.T) {
	mockRepo := newMockStatusRepo()
	mockRepo.statuses["0"] = &models.Status{ID: 1, Name: "In Stock"}
	mockRepo.statuses["1"] = &models.Status{ID: 2, Name: "Konsinyasi"}
	models.SetStatusRepository(mockRepo)
	app := fiber.New()
	app.Get("/statuses", GetStatuses)

	req := httptest.NewRequest("GET", "/statuses", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
