package utils

import (
	"errors"
	"unicode"

	"mori/pkg/models"

	"github.com/dchest/captcha" // Make sure you import the captcha package
)

// validate all fields when user registers, including captcha
func ValidateNewUser(user models.User, captchaID, captchaValue string) error {
	// 1) Validate Captcha first
	if err := validateCaptcha(captchaID, captchaValue); err != nil {
		return err
	}

	// 2) Then validate all user fields
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

// validateCaptcha calls dchest/captcha
func validateCaptcha(captchaID, captchaValue string) error {
	if captchaID == "" || captchaValue == "" {
		return errors.New("invalid captcha")
	}
	if !captcha.VerifyString(captchaID, captchaValue) {
		return errors.New("invalid captcha")
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

func validateEmail(email string) error {
	if fieldEmpty(email) {
		return errors.New("validation error")
	}
	return nil
}

// ---------------------------------------
// Manual password validation (no lookahead)
// ---------------------------------------
func validatePassword(password string) error {
	if fieldEmpty(password) {
		return errors.New("validation error")
	}

	// 1) At least 10 characters
	if len(password) < 10 {
		return errors.New("password must be at least 10 characters")
	}

	// 2) Must contain at least one uppercase, one digit, one special char
	var hasUpper, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case isSpecialChar(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// Helper to check if the character is "special" (not a letter/digit)
func isSpecialChar(ch rune) bool {
	// For a broad check, treat anything not a letter or digit as special
	return !(isLetter(ch) || isDigit(ch))
}

func fieldEmpty(value string) bool {
	return len(value) == 0
}

// If you need these:
func isLetter(ch rune) bool { return unicode.IsLetter(ch) }
func isDigit(ch rune) bool  { return unicode.IsDigit(ch) }
