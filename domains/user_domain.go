package domains

import (
	"errors"
	"strings"

	"gotask-api/datatransfers"
	"gotask-api/models"
	"gotask-api/repositories"
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
	// Session 13 adds bcrypt hashing right here
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // ⚠️ plain text — temporary!
	}

	return d.repo.Create(user)
}
