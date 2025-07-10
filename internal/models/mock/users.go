package mock

import "github.com/embracexyz/snippetbox/internal/models"

type MockUserModel struct {
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
