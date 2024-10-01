package accounts

import (
	"regexp"
)

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func ValidatePassword(password string) (bool, string) {

	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	digitRegex := regexp.MustCompile(`\d`)
	specialCharRegex := regexp.MustCompile(`[\W_]`)
	lengthRegex := regexp.MustCompile(`.{8,}`)

	if !lowercaseRegex.MatchString(password) {
		return true, "password must contain at least one lowercase letter"
	}
	if !uppercaseRegex.MatchString(password) {
		return true, "password must contain at least one uppercase letter"
	}
	if !digitRegex.MatchString(password) {
		return true, "password must contain at least one digit"
	}
	
	if !specialCharRegex.MatchString(password) {
		return true, "password must contain at least one special character"
	}
	if !lengthRegex.MatchString(password) {
		return true, "password must be at least 8 characters long"
	}
	return false, ""
}