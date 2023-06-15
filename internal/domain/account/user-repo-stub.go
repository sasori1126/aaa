package account

import "axis/ecommerce-backend/internal/models"

type UserRepoStub struct {
	users []models.User
}

func (s UserRepoStub) FindAll() ([]models.User, error) {
	return s.users, nil
}

func (s UserRepoStub) CreateUser() error {
	return nil
}

func NewUserRepoStub() UserRepoStub {
	users := []models.User{
		{
			Name: "Marvin",
			Email: "marvin@demo.com",
			PhoneNumber: "0704407117",
			IsActive: true,
		},
	}

	return UserRepoStub{users: users}
}
