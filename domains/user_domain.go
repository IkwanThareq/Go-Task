package domains

import (
	"errors"
	"strings"

	"gotask-api/datatransfers"
	"gotask-api/models"
	"gotask-api/repositories"
	"gotask-api/utils"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserDomain interface {
	Register(req datatransfers.RegisterRequest) (*models.User, error)
}

type userDomain struct {
	repo repositories.UserRepository
}

func NewUserDomain(repo repositories.UserRepository) UserDomain {
	return &userDomain{repo: repo}
}

func (d *userDomain) Register(req datatransfers.RegisterRequest) (*models.User, error) {
	// normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	// check if email already taken
	existing, _ := d.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// create user — password stored as plain text FOR NOW
	// Hashed password

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	}

	return d.repo.Create(user)
}
