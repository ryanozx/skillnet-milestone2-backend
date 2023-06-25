package database

import "github.com/ryanozx/skillnet-milestone2-backend/models"

type AuthDBHandler interface {
	GetUserByUsername(string) (*models.User, error)
}
