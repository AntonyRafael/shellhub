package services

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
)

type UserService interface {
	UpdateDataUser(ctx context.Context, user *models.User, id string) ([]string, error)
	UpdatePasswordUser(ctx context.Context, currentPassword, newPassword, id string) error
}

func (s *service) UpdateDataUser(ctx context.Context, user *models.User, id string) ([]string, error) {
	if _, _, err := s.store.UserGetByID(ctx, id, false); err != nil {
		return nil, NewErrUserNotFound(id, err)
	}

	if invalidFields, err := validator.ValidateStruct(user.UserData); err != nil {
		return invalidFields, NewErrUserInvalid(invalidFields, nil)
	}

	validator.FormatUser(user)

	var conflictFields []string
	var duplicatedValues []string
	existentUser, _ := s.store.UserGetByUsername(ctx, user.Username)
	if existentUser != nil && existentUser.ID != id {
		conflictFields = append(conflictFields, "username")
		duplicatedValues = append(duplicatedValues, user.Username)
	}

	existentUser, _ = s.store.UserGetByEmail(ctx, user.Email)
	if existentUser != nil && existentUser.ID != id {
		conflictFields = append(conflictFields, "email")
		duplicatedValues = append(duplicatedValues, user.Email)
	}

	if len(conflictFields) > 0 {
		return conflictFields, NewErrUserDuplicated(duplicatedValues, nil)
	}

	return nil, s.store.UserUpdateData(ctx, user, id)
}

func (s *service) UpdatePasswordUser(ctx context.Context, currentPassword, newPassword, id string) error {
	if !validator.ValidateFieldPassword(currentPassword) {
		return NewErrUserInvalid([]string{"current_password"}, nil)
	}

	if !validator.ValidateFieldPassword(newPassword) {
		return NewErrUserInvalid([]string{"new_password"}, nil)
	}

	currentPassword = validator.HashPassword(currentPassword)
	newPassword = validator.HashPassword(newPassword)
	if currentPassword == newPassword {
		return NewErrUserDuplicated([]string{"current_password", "new_password"}, nil)
	}

	user, _, err := s.store.UserGetByID(ctx, id, false)
	if user == nil {
		return NewErrUserNotFound(id, err)
	}

	if user.Password != currentPassword {
		return NewErrUserInvalid([]string{"current_password"}, nil)
	}

	return s.store.UserUpdatePassword(ctx, newPassword, id)
}
