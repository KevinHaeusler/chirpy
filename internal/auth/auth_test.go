package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestMakeAndValidateJWT_Success(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}

	if gotID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	rightSecret := "right-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, rightSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatalf("expected error for wrong secret, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		token := "Bearer valid-token"
		bearerToken, err := GetBearerToken(http.Header{"Authorization": []string{token}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bearerToken != "valid-token" {
			t.Fatalf("expected %q, got %q", "valid-token", bearerToken)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		bearerToken, err := GetBearerToken(http.Header{})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if bearerToken != "" {
			t.Fatalf("expected empty token, got %q", bearerToken)
		}
	})

	t.Run("no token after Bearer", func(t *testing.T) {
		bearerToken, err := GetBearerToken(http.Header{"Authorization": []string{"Bearer"}})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if bearerToken != "" {
			t.Fatalf("expected empty token, got %q", bearerToken)
		}
	})
}

func TestGetBearerToken_EmptyCases(t *testing.T) {
	tests := []struct {
		name   string
		header http.Header
	}{
		{
			name:   "missing header",
			header: http.Header{},
		},
		{
			name:   "empty value",
			header: http.Header{"Authorization": []string{""}},
		},
		{
			name:   "only Bearer",
			header: http.Header{"Authorization": []string{"Bearer"}},
		},
		{
			name:   "Bearer with spaces",
			header: http.Header{"Authorization": []string{"Bearer   "}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.header)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if token != "" {
				t.Fatalf("expected empty token, got %q", token)
			}
		})
	}
}
