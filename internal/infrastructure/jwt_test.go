package infrastructure_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure"
)

func TestJWTManager_GenerateAndValidate_Success(t *testing.T) {
	// Arrange
	secret := "test-secret"
	ttl := 1 * time.Minute
	userID := "user-123"

	jwtManager := infrastructure.NewJWTManager(secret, ttl)

	// Act
	token, err := jwtManager.Generate(userID)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	gotUserID, err := jwtManager.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	// Assert
	if gotUserID != userID {
		t.Errorf("expected userID %q, got %q", userID, gotUserID)
	}
}

func TestJWTManager_Validate_EmptyToken(t *testing.T) {
	jwtManager := infrastructure.NewJWTManager("secret", time.Minute)

	_, err := jwtManager.Validate("")

	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestJWTManager_Validate_InvalidToken(t *testing.T) {
	jwtManager := infrastructure.NewJWTManager("secret", time.Minute)

	_, err := jwtManager.Validate("this.is.not.a.jwt")

	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestJWTManager_Validate_WrongSecret(t *testing.T) {
	// Arrange
	managerA := infrastructure.NewJWTManager("secret-A", time.Minute)
	managerB := infrastructure.NewJWTManager("secret-B", time.Minute)

	token, err := managerA.Generate("user-123")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Act
	_, err = managerB.Validate(token)

	// Assert
	if err == nil {
		t.Fatal("expected error for token signed with wrong secret")
	}
}

func TestJWTManager_Validate_ExpiredToken(t *testing.T) {
	// Arrange
	jwtManager := infrastructure.NewJWTManager("secret", -1*time.Minute)

	token, err := jwtManager.Generate("user-123")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Act
	_, err = jwtManager.Validate(token)

	// Assert
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestJWTManager_Generate_TokenFormat(t *testing.T) {
	jwtManager := infrastructure.NewJWTManager("secret", time.Minute)

	token, err := jwtManager.Generate("user-123")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("expected JWT to have 3 parts, got %d", len(parts))
	}
}
