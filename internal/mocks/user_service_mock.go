// internal/mocks/user_service_mock.go
package mocks

import (
	"simdokpol/internal/models"

	"github.com/stretchr/testify/mock"
)

// UserService adalah mock sederhana untuk interface service pengguna.
// Nama harus persis "UserService" supaya test yang sekarang dapat compile.
type UserService struct {
	mock.Mock
}

// Create mock
func (m *UserService) Create(user *models.User, createdBy int) error {
	args := m.Called(user, createdBy)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

// UpdateProfile mock
func (m *UserService) UpdateProfile(userID int, user *models.User) (*models.User, error) {
	args := m.Called(userID, user)
	if u := args.Get(0); u != nil {
		if usr, ok := u.(*models.User); ok {
			return usr, args.Error(1)
		}
	}
	return nil, args.Error(1)
}
