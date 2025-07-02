package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/repo"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type UserService interface {
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
		currentPassword string,
		newPassword string,
	) (status_codes.UserUpdatePassword, error)
}

type userService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return userService{
		repo: repo,
	}
}

func (s userService) UpdateName(
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

func (s userService) UpdatePassword(
	ctx context.Context,
	user *entities.User,
	currentPassword string,
	newPassword string,
) (status_codes.UserUpdatePassword, error) {
	// Check if the current password matches the db one
	if !util.CheckPasswordHash(currentPassword, user.Password) {
		return status_codes.UserUpdatePasswordWrongCurrent, nil
	}

	// Validate the new password
	if !rules.IsValidUserPassword(newPassword) {
		return status_codes.UserUpdatePasswordInvalid, nil
	}

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
