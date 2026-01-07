// form/validate_test.go
package form

import "testing"

func TestRequired(t *testing.T) {
	v := Required()

	if v("") == nil {
		t.Error("empty string should fail Required")
	}
	if v("hello") != nil {
		t.Error("non-empty string should pass Required")
	}
	if v(nil) == nil {
		t.Error("nil should fail Required")
	}
}

func TestMinLength(t *testing.T) {
	v := MinLength(3)

	if v("ab") == nil {
		t.Error("'ab' should fail MinLength(3)")
	}
	if v("abc") != nil {
		t.Error("'abc' should pass MinLength(3)")
	}
	if v("abcd") != nil {
		t.Error("'abcd' should pass MinLength(3)")
	}
}

func TestMaxLength(t *testing.T) {
	v := MaxLength(5)

	if v("hello") != nil {
		t.Error("'hello' should pass MaxLength(5)")
	}
	if v("hello!") == nil {
		t.Error("'hello!' should fail MaxLength(5)")
	}
}

func TestPattern(t *testing.T) {
	v := Pattern(`^[a-z]+$`, "lowercase only")

	if v("hello") != nil {
		t.Error("'hello' should pass pattern")
	}
	if v("Hello") == nil {
		t.Error("'Hello' should fail pattern")
	}
	if v("123") == nil {
		t.Error("'123' should fail pattern")
	}
}

func TestEmail(t *testing.T) {
	v := Email()

	if v("test@example.com") != nil {
		t.Error("valid email should pass")
	}
	if v("not-an-email") == nil {
		t.Error("invalid email should fail")
	}
	if v("@example.com") == nil {
		t.Error("missing local part should fail")
	}
}

func TestMinSelected(t *testing.T) {
	v := MinSelected(2)

	if v([]string{"a"}) == nil {
		t.Error("1 item should fail MinSelected(2)")
	}
	if v([]string{"a", "b"}) != nil {
		t.Error("2 items should pass MinSelected(2)")
	}
}

func TestMaxSelected(t *testing.T) {
	v := MaxSelected(2)

	if v([]string{"a", "b"}) != nil {
		t.Error("2 items should pass MaxSelected(2)")
	}
	if v([]string{"a", "b", "c"}) == nil {
		t.Error("3 items should fail MaxSelected(2)")
	}
}

func TestComposeValidators(t *testing.T) {
	validators := []Validator{Required(), MinLength(3)}

	// Empty fails Required first
	err := Compose(validators...)("")
	if err == nil {
		t.Error("should fail")
	}

	// "ab" passes Required, fails MinLength
	err = Compose(validators...)("ab")
	if err == nil {
		t.Error("should fail MinLength")
	}

	// "abc" passes both
	err = Compose(validators...)("abc")
	if err != nil {
		t.Error("should pass both validators")
	}
}

func TestMinLengthNonString(t *testing.T) {
	v := MinLength(3)

	// Non-string types should pass (return nil)
	if v(123) != nil {
		t.Error("non-string should pass MinLength")
	}
	if v(true) != nil {
		t.Error("bool should pass MinLength")
	}
}

func TestMaxLengthNonString(t *testing.T) {
	v := MaxLength(5)

	// Non-string types should pass (return nil)
	if v(123) != nil {
		t.Error("non-string should pass MaxLength")
	}
	if v(true) != nil {
		t.Error("bool should pass MaxLength")
	}
}

func TestPatternNonString(t *testing.T) {
	v := Pattern(`^[a-z]+$`, "lowercase only")

	// Non-string types should pass (return nil)
	if v(123) != nil {
		t.Error("non-string should pass Pattern")
	}
	if v(true) != nil {
		t.Error("bool should pass Pattern")
	}
}

func TestMinSelectedNonArray(t *testing.T) {
	v := MinSelected(2)

	// Non-array types should pass (return nil)
	if v("hello") != nil {
		t.Error("non-array should pass MinSelected")
	}
	if v(123) != nil {
		t.Error("int should pass MinSelected")
	}
}

func TestMaxSelectedNonArray(t *testing.T) {
	v := MaxSelected(2)

	// Non-array types should pass (return nil)
	if v("hello") != nil {
		t.Error("non-array should pass MaxSelected")
	}
	if v(123) != nil {
		t.Error("int should pass MaxSelected")
	}
}
