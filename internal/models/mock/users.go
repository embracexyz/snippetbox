package mock

import (
	"time"

	"github.com/embracexyz/snippetbox/internal/models"
)

type MockUserModel struct {
}

func (m *MockUserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword == "pa$$word" {
			return nil
		}
		return models.ErrInvalidCredentials
	}
	return models.ErrNoRecord
}

func (m *MockUserModel) GetUser(id int) (*models.User, error) {
	switch id {
	case 1:
		return &models.User{
			ID:      1,
			Name:    "Alice",
			Email:   "alice@example.com",
			Created: time.Now(),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *MockUserModel) Insert(name, email, password string) error {
	switch email {
	case "alice@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}

}

func (m *MockUserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *MockUserModel) Exists(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, nil
}
