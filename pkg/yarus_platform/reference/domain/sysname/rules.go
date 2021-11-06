package sysname

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	Regexp = "^[a-zA-Z0-9_\\.]+$"
)

var ValidationRules = []validation.Rule{
	validation.Required,
	validation.Length(2, 100),
	validation.Match(regexp.MustCompile(Regexp)),
}
