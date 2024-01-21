package validation

import (
	"regexp"
)

type Validator struct {
	EmailRegex *regexp.Regexp
}

func New() *Validator {
	regex := regexp.MustCompile(`[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*`)
	return &Validator{EmailRegex: regex}
}

func (v *Validator) IsEmail(email string) bool {
	return v.EmailRegex.MatchString(email)
}
