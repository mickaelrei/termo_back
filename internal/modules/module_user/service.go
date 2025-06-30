package module_user

import (
	"context"
	"fmt"
	"log"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type Service interface {
	// UpdateName ensures the new name is valid and then changes it
	UpdateName(
		ctx context.Context,
		user *entities.User,
		newName string,
	) (status_codes.UserUpdateName, error)

	// UpdatePassword ensures the new password is valid and then changes it
	UpdatePassword(
		ctx context.Context,
		user *entities.User,
		newPassword string,
	) (status_codes.UserUpdatePassword, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) UpdateName(
	ctx context.Context,
	user *entities.User,
	newName string,
) (status_codes.UserUpdateName, error) {
	// Clean name and validate
	newName = strings.TrimSpace(newName)
	if !rules.IsValidUserName(newName) {
		return status_codes.UserUpdateNameInvalid, nil
	}

	err := s.repo.UpdateName(ctx, user.ID, newName)
	if err != nil {
		log.Printf("[UpdateName] | %v", err)
		return -1, err
	}

	return status_codes.UserUpdateNameSuccess, nil
}

func (s service) UpdatePassword(
	ctx context.Context,
	user *entities.User,
	newPassword string,
) (status_codes.UserUpdatePassword, error) {
	// Validate password
	if !rules.IsValidUserPassword(newPassword) {
		return status_codes.UserUpdatePasswordInvalid, nil
	}

	// Hash the password
	newPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return -1, fmt.Errorf("[HashPassword] | %v", err)
	}

	err = s.repo.UpdatePassword(ctx, user.ID, newPassword)
	if err != nil {
		return -1, fmt.Errorf("[UpdatePassword] | %v", err)
	}

	return status_codes.UserUpdatePasswordSuccess, nil
}
