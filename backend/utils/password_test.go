package utils

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "Secure@123", false},
		{"too short", "Sec@1", true},
		{"no uppercase", "secure@123", true},
		{"no lowercase", "SECURE@123", true},
		{"no digit", "Secure@abc", true},
		{"no special char", "Secure1234", true},
		{"exactly 8 chars valid", "Abc@1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "Secure@123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == password {
		t.Fatal("HashPassword() returned plaintext password unchanged")
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("CheckPasswordHash() should succeed for correct password, got: %v", err)
	}

	if err := CheckPasswordHash("WrongPass@1", hash); err == nil {
		t.Error("CheckPasswordHash() should fail for wrong password")
	}
}

func TestHashPasswordIsDeterministicallyDifferent(t *testing.T) {
	// bcrypt should produce different hashes on each call (salted)
	password := "Secure@123"
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)
	if hash1 == hash2 {
		t.Error("HashPassword() produced identical hashes for same input — bcrypt salt not working")
	}
}
