// form/validate.go
package shell

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validator is a function that validates a value.
type Validator func(value any) error

// Required validates that a value is not empty.
func Required() Validator {
	return func(value any) error {
		if value == nil {
			return errors.New("required")
		}
		if s, ok := value.(string); ok && strings.TrimSpace(s) == "" {
			return errors.New("required")
		}
		return nil
	}
}

// MinLength validates minimum string length.
func MinLength(n int) Validator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		if len(s) < n {
			return fmt.Errorf("minimum %d characters required", n)
		}
		return nil
	}
}

// MaxLength validates maximum string length.
func MaxLength(n int) Validator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		if len(s) > n {
			return fmt.Errorf("maximum %d characters allowed", n)
		}
		return nil
	}
}

// Pattern validates against a regex pattern.
func Pattern(pattern string, msg string) Validator {
	re := regexp.MustCompile(pattern)
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		if !re.MatchString(s) {
			return errors.New(msg)
		}
		return nil
	}
}

// Email validates email format.
func Email() Validator {
	return Pattern(`^[^@\s]+@[^@\s]+\.[^@\s]+$`, "invalid email format")
}

// MinSelected validates minimum selections for multi-select.
func MinSelected(n int) Validator {
	return func(value any) error {
		if arr, ok := value.([]string); ok {
			if len(arr) < n {
				return fmt.Errorf("select at least %d options", n)
			}
		}
		return nil
	}
}

// MaxSelected validates maximum selections for multi-select.
func MaxSelected(n int) Validator {
	return func(value any) error {
		if arr, ok := value.([]string); ok {
			if len(arr) > n {
				return fmt.Errorf("select at most %d options", n)
			}
		}
		return nil
	}
}

// Compose combines multiple validators into one.
func Compose(validators ...Validator) Validator {
	return func(value any) error {
		for _, v := range validators {
			if err := v(value); err != nil {
				return err
			}
		}
		return nil
	}
}
