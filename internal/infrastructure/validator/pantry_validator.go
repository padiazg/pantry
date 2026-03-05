package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// PantryValidator Validador central para entidades de la despensa: EAN-13, stock levels y movimientos
type PantryValidator struct {
	errors []string
}

// NewPantryValidator creates a new PantryValidator instance
func NewPantryValidator() *PantryValidator {
	return &PantryValidator{
		errors: make([]string, 0),
	}
}

// Validate runs all validation rules
func (v *PantryValidator) Validate(data interface{}) error {
	v.errors = make([]string, 0)

	// TODO: Implement validation logic
	// Example:
	// if user, ok := data.(*User); ok {
	//     v.validateEmail(user.Email)
	//     v.validateName(user.Name)
	// }

	if len(v.errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(v.errors, ", "))
	}

	return nil
}

// Required validates that a field is not empty
func (v *PantryValidator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, fmt.Sprintf("%s is required", field))
	}
}

// Email validates an email address
func (v *PantryValidator) Email(field, value string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors = append(v.errors, fmt.Sprintf("%s must be a valid email", field))
	}
}

// MinLength validates minimum string length
func (v *PantryValidator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
}

// MaxLength validates maximum string length
func (v *PantryValidator) MaxLength(field, value string, max int) {
	if len(value) > max {
		v.errors = append(v.errors, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
}

// Range validates numeric range
func (v *PantryValidator) Range(field string, value, min, max int) {
	if value < min || value > max {
		v.errors = append(v.errors, fmt.Sprintf("%s must be between %d and %d", field, min, max))
	}
}

// HasErrors returns true if there are validation errors
func (v *PantryValidator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors
func (v *PantryValidator) Errors() []string {
	return v.errors
}
