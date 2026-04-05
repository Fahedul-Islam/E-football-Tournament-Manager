package user

import (
	"context"
	"errors"
	"testing"
	"tournament-manager/config"
	"tournament-manager/internal/domain"
)

// mockUserRepo is a minimal in-memory stub for testing the service layer.
type mockUserRepo struct {
	registered []domain.User
	shouldFail bool
}

func (m *mockUserRepo) Register(ctx context.Context, user domain.User) error {
	if m.shouldFail {
		return errors.New("db error")
	}
	m.registered = append(m.registered, user)
	return nil
}

func (m *mockUserRepo) GetUserData(ctx context.Context, email, password, role string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, id int, user domain.User) error {
	return errors.New("not implemented")
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id int) error {
	return errors.New("not implemented")
}

func newTestService(repo *mockUserRepo) Service {
	cfg := &config.Config{}
	cfg.JWT.Secret = []byte("test-secret")
	return NewUserService(cfg, repo)
}

func TestRegister_WeakPasswordRejected(t *testing.T) {
	repo := &mockUserRepo{}
	svc := newTestService(repo)

	// Weak passwords must be rejected before anything reaches the repo
	weakPasswords := []string{
		"short",        // too short
		"alllowercase1!", // no uppercase
		"ALLUPPERCASE1!", // no lowercase
		"NoDigits!abc",   // no digit
		"NoSpecial123",   // no special char
	}

	for _, pw := range weakPasswords {
		err := svc.Register(context.Background(), domain.User{
			Email:        "test@example.com",
			PasswordHash: pw,
			Username:     "testuser",
			Role:         "player",
		})
		if err == nil {
			t.Errorf("Register() should reject weak password %q but did not", pw)
		}
		if len(repo.registered) > 0 {
			t.Errorf("Register() with weak password %q should not reach the repo", pw)
		}
	}
}

func TestRegister_StrongPasswordAccepted(t *testing.T) {
	repo := &mockUserRepo{}
	svc := newTestService(repo)

	err := svc.Register(context.Background(), domain.User{
		Email:        "test@example.com",
		PasswordHash: "Secure@123",
		Username:     "testuser",
		Role:         "player",
	})
	if err != nil {
		t.Fatalf("Register() should succeed with strong password, got: %v", err)
	}
	if len(repo.registered) != 1 {
		t.Fatal("Register() should have persisted user to repo")
	}
	stored := repo.registered[0]
	if stored.PasswordHash == "Secure@123" {
		t.Error("Register() must hash the password before storing — plaintext detected in repo")
	}
	if len(stored.PasswordHash) < 30 {
		t.Error("Register() stored password looks too short to be a bcrypt hash")
	}
}
