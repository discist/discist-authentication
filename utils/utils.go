package utils

import (
	"strings"
)

func NormalizeEmail(email string) string {

	nospace := strings.Replace(email, " ", "", -1)
	return strings.TrimSpace(strings.ToLower(nospace))

}
