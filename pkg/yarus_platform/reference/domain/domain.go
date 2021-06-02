package domain

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	SysnameRegexp = "^[a-z0-9_]+$"
)

var SysnameValidationRules = []validation.Rule{
	validation.Required,
	validation.Length(2, 100),
	validation.Match(regexp.MustCompile(SysnameRegexp)),
}
