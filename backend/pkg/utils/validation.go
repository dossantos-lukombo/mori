package utils

import (
	"errors"
	"mori/pkg/models"
)

// validate all fields when user registers
func ValidateNewUser(user models.User) error {
	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}
	if err := validateLastName(user.LastName); err != nil {
		return err
	}
	if err := validateBirth(user.DateOfBirth); err != nil {
		return err
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	return nil
}
func validateFirstName(name string) error {
	if fieldEmpty(name) {
		return errors.New("validation error")
	}
	return nil
}
func validateLastName(name string) error {
	if fieldEmpty(name) {
		return errors.New("validation error")
	}
	return nil
}
func validateBirth(birthDate string) error {
	if fieldEmpty(birthDate) {
		return errors.New("validation error")
	}
	return nil
}
func validatePassword(password string) error {
	if fieldEmpty(password) {
		return errors.New("validation error")
	}
	return nil
}
func validateEmail(email string) error {
	if fieldEmpty(email) {
		return errors.New("validation error")
	}
	return nil
}

func fieldEmpty(value string) bool {
	return len(value) == 0
}
