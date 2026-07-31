package validator

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	emailRegex     *regexp.Regexp
	emailRegexOnce sync.Once
	passwordRegex  *regexp.Regexp
	passwordOnce   sync.Once
)

func getEmailRegex() *regexp.Regexp {
	emailRegexOnce.Do(func() {
		emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	})
	return emailRegex
}

func getPasswordRegex() *regexp.Regexp {
	passwordOnce.Do(func() {
		passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	})
	return passwordRegex
}

func RegisterCustomValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	if err := v.RegisterValidation("email_format", validateEmailFormat); err != nil {
		return err
	}

	if err := v.RegisterValidation("password_strong", validatePasswordStrong); err != nil {
		return err
	}

	return nil
}

func validateEmailFormat(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true
	}
	return getEmailRegex().MatchString(email)
}

func validatePasswordStrong(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true
	}
	return getPasswordRegex().MatchString(password)
}

func IsValidEmail(email string) bool {
	return getEmailRegex().MatchString(email)
}

func IsStrongPassword(password string) bool {
	return getPasswordRegex().MatchString(password)
}

func GetPasswordStrengthMessage() string {
	return "Password must contain at least 8 characters, including uppercase, lowercase, digit, and special character"
}
