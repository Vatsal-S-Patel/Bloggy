package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fieldLevel validator.FieldLevel) bool {
	re := `.*[!@#\$%\^&\*\(\)_\+\-=\[\]\{\}\\|;:'",<>\.\?/` + "`" + `~].*`
	return regexp.MustCompile(re).MatchString(fieldLevel.Field().String())
}
