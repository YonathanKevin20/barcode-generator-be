// utils/token_test.go
package utils

import (
	"barcode-generator-be/models"
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	user := &models.User{ID: 1, Username: "testuser", Role: "admin"}
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected token, got empty string")
	}
}
