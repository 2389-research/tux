# Forms Package Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a language-agnostic form system with 7 field types, validation, and modal integration.

**Architecture:** New `form` package with Field interface, concrete field types backed by huh, Form compositor with group support, and FormModal adapter in modal package.

**Tech Stack:** Go, charmbracelet/huh, charmbracelet/bubbletea, charmbracelet/lipgloss

---

### Task 1: Field Interface and Base Types

**Files:**
- Create: `form/field.go`
- Create: `form/form.go`
- Create: `form/field_test.go`

**Step 1: Write the failing test**

```go
// form/field_test.go
package form

import "testing"

func TestFieldInterface(t *testing.T) {
	// Verify Option helper works
	opt := Option("Display", "value")
	if opt.Label != "Display" {
		t.Errorf("expected label 'Display', got %s", opt.Label)
	}
	if opt.Value != "value" {
		t.Errorf("expected value 'value', got %v", opt.Value)
	}
}

func TestValuesAccessors(t *testing.T) {
	v := Values{
		"name":     "Alice",
		"active":   true,
		"features": []string{"a", "b"},
	}

	if v.String("name") != "Alice" {
		t.Errorf("expected 'Alice', got %s", v.String("name"))
	}
	if v.String("missing") != "" {
		t.Error("missing key should return empty string")
	}
	if !v.Bool("active") {
		t.Error("expected true")
	}
	if v.Bool("missing") {
		t.Error("missing bool should return false")
	}
	strs := v.Strings("features")
	if len(strs) != 2 || strs[0] != "a" {
		t.Errorf("expected [a, b], got %v", strs)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v`
Expected: FAIL - package doesn't exist

**Step 3: Write minimal implementation**

```go
// form/field.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

// Field is the interface for form fields.
type Field interface {
	// Identity
	ID() string
	Label() string

	// State
	Value() any
	SetValue(any)
	Focused() bool
	Focus()
	Blur()

	// Validation
	Validate() error

	// Rendering
	Render(width int, th theme.Theme, focused bool) string

	// Input handling
	HandleKey(key tea.KeyMsg) bool

	// Lifecycle
	Init() tea.Cmd
}

// SelectOption represents an option in a select field.
type SelectOption[T any] struct {
	Label string
	Value T
}

// Option creates a SelectOption.
func Option[T any](label string, value T) SelectOption[T] {
	return SelectOption[T]{Label: label, Value: value}
}

// Values holds form field values by ID.
type Values map[string]any

// String returns a string value or empty string if not found/wrong type.
func (v Values) String(id string) string {
	if val, ok := v[id].(string); ok {
		return val
	}
	return ""
}

// Bool returns a bool value or false if not found/wrong type.
func (v Values) Bool(id string) bool {
	if val, ok := v[id].(bool); ok {
		return val
	}
	return false
}

// Strings returns a []string value or nil if not found/wrong type.
func (v Values) Strings(id string) []string {
	if val, ok := v[id].([]string); ok {
		return val
	}
	return nil
}

// Int returns an int value or 0 if not found/wrong type.
func (v Values) Int(id string) int {
	if val, ok := v[id].(int); ok {
		return val
	}
	return 0
}
```

```go
// form/form.go
package form

// State represents the form state.
type State int

const (
	StateActive State = iota
	StateSubmitted
	StateCancelled
)
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/
git commit -m "feat(form): add Field interface and Values type"
```

---

### Task 2: Validators

**Files:**
- Create: `form/validate.go`
- Create: `form/validate_test.go`

**Step 1: Write the failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestRequired`
Expected: FAIL - Validator not defined

**Step 3: Write minimal implementation**

```go
// form/validate.go
package form

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/validate.go form/validate_test.go
git commit -m "feat(form): add validators"
```

---

### Task 3: Huh Theme Bridge

**Files:**
- Create: `form/huh_theme.go`
- Create: `form/huh_theme_test.go`

**Step 1: Write the failing test**

```go
// form/huh_theme_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestToHuhTheme(t *testing.T) {
	th := theme.Get("dracula")
	huhTheme := ToHuhTheme(th)

	if huhTheme == nil {
		t.Fatal("expected non-nil huh theme")
	}

	// Verify focused title uses primary color
	// (huh theme is configured, just verify it returns something)
}

func TestToHuhThemeAllThemes(t *testing.T) {
	themes := []string{"dracula", "nord", "gruvbox", "high-contrast", "neo-terminal"}

	for _, name := range themes {
		t.Run(name, func(t *testing.T) {
			th := theme.Get(name)
			huhTheme := ToHuhTheme(th)
			if huhTheme == nil {
				t.Errorf("ToHuhTheme returned nil for %s", name)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestToHuhTheme`
Expected: FAIL - ToHuhTheme not defined

**Step 3: Write minimal implementation**

```go
// form/huh_theme.go
package form

import (
	"github.com/2389-research/tux/theme"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ToHuhTheme converts a tux theme to a huh theme.
func ToHuhTheme(th theme.Theme) *huh.Theme {
	t := huh.ThemeBase()

	// Focused state styles
	t.Focused.Base = t.Focused.Base.
		BorderForeground(th.Primary())

	t.Focused.Title = lipgloss.NewStyle().
		Foreground(th.Primary()).
		Bold(true)

	t.Focused.Description = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(th.Primary())

	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(th.Success()).
		Bold(true)

	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(th.Background()).
		Background(th.Primary()).
		Bold(true).
		Padding(0, 1)

	t.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(th.Foreground()).
		Background(th.Border()).
		Padding(0, 1)

	t.Focused.TextInput.Cursor = lipgloss.NewStyle().
		Foreground(th.Primary())

	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Focused.TextInput.Prompt = lipgloss.NewStyle().
		Foreground(th.Primary())

	// Blurred state - more muted
	t.Blurred.Base = t.Blurred.Base.
		BorderForeground(th.Border())

	t.Blurred.Title = lipgloss.NewStyle().
		Foreground(th.Foreground())

	t.Blurred.Description = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Blurred.SelectSelector = lipgloss.NewStyle().
		Foreground(th.Muted())

	t.Blurred.SelectedOption = lipgloss.NewStyle().
		Foreground(th.Foreground())

	t.Blurred.UnselectedOption = lipgloss.NewStyle().
		Foreground(th.Muted())

	// Error styling
	t.Focused.ErrorMessage = lipgloss.NewStyle().
		Foreground(th.Error())

	t.Focused.ErrorIndicator = lipgloss.NewStyle().
		Foreground(th.Error())

	return t
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/huh_theme.go form/huh_theme_test.go
git commit -m "feat(form): add huh theme bridge"
```

---

### Task 4: Input Field

**Files:**
- Create: `form/input.go`
- Create: `form/input_test.go`

**Step 1: Write the failing test**

```go
// form/input_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestInputField(t *testing.T) {
	f := NewInput().
		ID("username").
		Label("Username").
		Placeholder("Enter username...")

	if f.ID() != "username" {
		t.Errorf("expected ID 'username', got %s", f.ID())
	}
	if f.Label() != "Username" {
		t.Errorf("expected label 'Username', got %s", f.Label())
	}
}

func TestInputFieldValue(t *testing.T) {
	f := NewInput().ID("test")

	f.SetValue("hello")
	if f.Value().(string) != "hello" {
		t.Errorf("expected 'hello', got %v", f.Value())
	}
}

func TestInputFieldValidation(t *testing.T) {
	f := NewInput().
		ID("email").
		Validate(Required(), Email())

	f.SetValue("")
	if f.Validate() == nil {
		t.Error("empty should fail validation")
	}

	f.SetValue("not-email")
	if f.Validate() == nil {
		t.Error("invalid email should fail")
	}

	f.SetValue("test@example.com")
	if f.Validate() != nil {
		t.Error("valid email should pass")
	}
}

func TestInputFieldRender(t *testing.T) {
	f := NewInput().
		ID("test").
		Label("Test Input")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}

func TestInputFieldFocus(t *testing.T) {
	f := NewInput().ID("test")

	if f.Focused() {
		t.Error("should not be focused initially")
	}

	f.Focus()
	if !f.Focused() {
		t.Error("should be focused after Focus()")
	}

	f.Blur()
	if f.Focused() {
		t.Error("should not be focused after Blur()")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestInputField`
Expected: FAIL - NewInput not defined

**Step 3: Write minimal implementation**

```go
// form/input.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// InputField is a single-line text input.
type InputField struct {
	id          string
	label       string
	placeholder string
	value       string
	validators  []Validator
	focused     bool
	huhField    *huh.Input
	huhTheme    *huh.Theme
}

// NewInput creates a new input field.
func NewInput() *InputField {
	return &InputField{}
}

// ID sets the field ID.
func (f *InputField) ID(id string) *InputField {
	f.id = id
	return f
}

// Label sets the field label.
func (f *InputField) Label(label string) *InputField {
	f.label = label
	return f
}

// Placeholder sets the placeholder text.
func (f *InputField) Placeholder(placeholder string) *InputField {
	f.placeholder = placeholder
	return f
}

// Validate adds validators to the field.
func (f *InputField) Validate(validators ...Validator) *InputField {
	f.validators = append(f.validators, validators...)
	return f
}

// Field interface implementation

func (f *InputField) id_() string    { return f.id }
func (f *InputField) label_() string { return f.label }

// Implement the Field interface methods with proper names
func (f *InputField) GetID() string    { return f.id }
func (f *InputField) GetLabel() string { return f.label }

func (f *InputField) Value() any {
	return f.value
}

func (f *InputField) SetValue(v any) {
	if s, ok := v.(string); ok {
		f.value = s
	}
}

func (f *InputField) Focused() bool {
	return f.focused
}

func (f *InputField) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

func (f *InputField) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

func (f *InputField) Validate_() error {
	return Compose(f.validators...)(f.value)
}

func (f *InputField) Init() tea.Cmd {
	f.huhField = huh.NewInput().
		Title(f.label).
		Placeholder(f.placeholder).
		Value(&f.value)

	if f.huhTheme != nil {
		// Theme is applied at form level
	}

	return f.huhField.Init()
}

func (f *InputField) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

func (f *InputField) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}

// Ensure InputField implements a simpler interface for now
// We'll reconcile with Field interface when building Form

func (f *InputField) ID_Get() string       { return f.id }
func (f *InputField) Label_Get() string    { return f.label }
func (f *InputField) Validate() error      { return Compose(f.validators...)(f.value) }
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/input.go form/input_test.go
git commit -m "feat(form): add Input field"
```

---

### Task 5: Select Field

**Files:**
- Create: `form/select.go`
- Create: `form/select_test.go`

**Step 1: Write the failing test**

```go
// form/select_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestSelectField(t *testing.T) {
	f := NewSelect[string]().
		ID("theme").
		Label("Choose Theme").
		Options(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		)

	if f.ID_Get() != "theme" {
		t.Errorf("expected ID 'theme', got %s", f.ID_Get())
	}
}

func TestSelectFieldDefault(t *testing.T) {
	f := NewSelect[string]().
		ID("theme").
		Options(
			Option("Dracula", "dracula"),
			Option("Nord", "nord"),
		).
		Default("nord")

	if f.Value().(string) != "nord" {
		t.Errorf("expected default 'nord', got %v", f.Value())
	}
}

func TestSelectFieldValue(t *testing.T) {
	f := NewSelect[string]().ID("test").Options(
		Option("A", "a"),
		Option("B", "b"),
	)

	f.SetValue("b")
	if f.Value().(string) != "b" {
		t.Errorf("expected 'b', got %v", f.Value())
	}
}

func TestSelectFieldRender(t *testing.T) {
	f := NewSelect[string]().
		ID("test").
		Label("Pick One").
		Options(Option("A", "a"))

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestSelectField`
Expected: FAIL - NewSelect not defined

**Step 3: Write minimal implementation**

```go
// form/select.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// SelectField is a single-choice selection field.
type SelectField[T comparable] struct {
	id         string
	label      string
	options    []SelectOption[T]
	value      T
	validators []Validator
	focused    bool
	huhField   *huh.Select[T]
	huhTheme   *huh.Theme
}

// NewSelect creates a new select field.
func NewSelect[T comparable]() *SelectField[T] {
	return &SelectField[T]{}
}

// ID sets the field ID.
func (f *SelectField[T]) ID(id string) *SelectField[T] {
	f.id = id
	return f
}

// Label sets the field label.
func (f *SelectField[T]) Label(label string) *SelectField[T] {
	f.label = label
	return f
}

// Options sets the available options.
func (f *SelectField[T]) Options(options ...SelectOption[T]) *SelectField[T] {
	f.options = options
	return f
}

// Default sets the default selected value.
func (f *SelectField[T]) Default(value T) *SelectField[T] {
	f.value = value
	return f
}

// Validate adds validators.
func (f *SelectField[T]) Validate(validators ...Validator) *SelectField[T] {
	f.validators = append(f.validators, validators...)
	return f
}

// Field interface implementation

func (f *SelectField[T]) ID_Get() string    { return f.id }
func (f *SelectField[T]) Label_Get() string { return f.label }

func (f *SelectField[T]) Value() any {
	return f.value
}

func (f *SelectField[T]) SetValue(v any) {
	if val, ok := v.(T); ok {
		f.value = val
	}
}

func (f *SelectField[T]) Focused() bool {
	return f.focused
}

func (f *SelectField[T]) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

func (f *SelectField[T]) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

func (f *SelectField[T]) Validate_() error {
	return Compose(f.validators...)(f.value)
}

func (f *SelectField[T]) Init() tea.Cmd {
	huhOpts := make([]huh.Option[T], len(f.options))
	for i, opt := range f.options {
		huhOpts[i] = huh.NewOption(opt.Label, opt.Value)
	}

	f.huhField = huh.NewSelect[T]().
		Title(f.label).
		Options(huhOpts...).
		Value(&f.value)

	return f.huhField.Init()
}

func (f *SelectField[T]) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

func (f *SelectField[T]) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/select.go form/select_test.go
git commit -m "feat(form): add Select field"
```

---

### Task 6: Confirm Field

**Files:**
- Create: `form/confirm.go`
- Create: `form/confirm_test.go`

**Step 1: Write the failing test**

```go
// form/confirm_test.go
package form

import (
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestConfirmField(t *testing.T) {
	f := NewConfirm().
		ID("delete").
		Label("Delete this file?").
		Affirmative("Yes, delete").
		Negative("Cancel")

	if f.ID_Get() != "delete" {
		t.Errorf("expected ID 'delete', got %s", f.ID_Get())
	}
}

func TestConfirmFieldValue(t *testing.T) {
	f := NewConfirm().ID("test")

	// Default should be false
	if f.Value().(bool) != false {
		t.Error("default should be false")
	}

	f.SetValue(true)
	if f.Value().(bool) != true {
		t.Error("should be true after SetValue(true)")
	}
}

func TestConfirmFieldRender(t *testing.T) {
	f := NewConfirm().
		ID("test").
		Label("Confirm?")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestConfirmField`
Expected: FAIL - NewConfirm not defined

**Step 3: Write minimal implementation**

```go
// form/confirm.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/2389-research/tux/theme"
)

// ConfirmField is a yes/no confirmation field.
type ConfirmField struct {
	id          string
	label       string
	affirmative string
	negative    string
	value       bool
	focused     bool
	huhField    *huh.Confirm
	huhTheme    *huh.Theme
}

// NewConfirm creates a new confirm field.
func NewConfirm() *ConfirmField {
	return &ConfirmField{
		affirmative: "Yes",
		negative:    "No",
	}
}

// ID sets the field ID.
func (f *ConfirmField) ID(id string) *ConfirmField {
	f.id = id
	return f
}

// Label sets the field label.
func (f *ConfirmField) Label(label string) *ConfirmField {
	f.label = label
	return f
}

// Affirmative sets the affirmative button text.
func (f *ConfirmField) Affirmative(text string) *ConfirmField {
	f.affirmative = text
	return f
}

// Negative sets the negative button text.
func (f *ConfirmField) Negative(text string) *ConfirmField {
	f.negative = text
	return f
}

// Field interface implementation

func (f *ConfirmField) ID_Get() string    { return f.id }
func (f *ConfirmField) Label_Get() string { return f.label }

func (f *ConfirmField) Value() any {
	return f.value
}

func (f *ConfirmField) SetValue(v any) {
	if b, ok := v.(bool); ok {
		f.value = b
	}
}

func (f *ConfirmField) Focused() bool {
	return f.focused
}

func (f *ConfirmField) Focus() {
	f.focused = true
	if f.huhField != nil {
		f.huhField.Focus()
	}
}

func (f *ConfirmField) Blur() {
	f.focused = false
	if f.huhField != nil {
		f.huhField.Blur()
	}
}

func (f *ConfirmField) Validate() error {
	return nil // Confirm has no validation
}

func (f *ConfirmField) Init() tea.Cmd {
	f.huhField = huh.NewConfirm().
		Title(f.label).
		Affirmative(f.affirmative).
		Negative(f.negative).
		Value(&f.value)

	return f.huhField.Init()
}

func (f *ConfirmField) HandleKey(key tea.KeyMsg) bool {
	if f.huhField == nil {
		return false
	}
	_, cmd := f.huhField.Update(key)
	return cmd != nil
}

func (f *ConfirmField) Render(width int, th theme.Theme, focused bool) string {
	if f.huhField == nil {
		f.huhTheme = ToHuhTheme(th)
		f.Init()
	}
	return f.huhField.View()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/confirm.go form/confirm_test.go
git commit -m "feat(form): add Confirm field"
```

---

### Task 7: Note Field (Display Only)

**Files:**
- Create: `form/note.go`
- Create: `form/note_test.go`

**Step 1: Write the failing test**

```go
// form/note_test.go
package form

import (
	"strings"
	"testing"

	"github.com/2389-research/tux/theme"
)

func TestNoteField(t *testing.T) {
	f := NewNote().
		ID("warning").
		Title("Warning").
		Content("This action cannot be undone.")

	if f.ID_Get() != "warning" {
		t.Errorf("expected ID 'warning', got %s", f.ID_Get())
	}
}

func TestNoteFieldValue(t *testing.T) {
	f := NewNote().ID("test").Content("hello")

	// Note value is always nil (display only)
	if f.Value() != nil {
		t.Error("note value should be nil")
	}
}

func TestNoteFieldRender(t *testing.T) {
	f := NewNote().
		ID("test").
		Title("Important").
		Content("Read this carefully.")

	th := theme.Get("dracula")
	output := f.Render(40, th, true)

	if output == "" {
		t.Error("render should produce output")
	}
	if !strings.Contains(output, "Important") && !strings.Contains(output, "Read this") {
		t.Error("render should contain title or content")
	}
}

func TestNoteFieldNoValidation(t *testing.T) {
	f := NewNote().ID("test")

	if f.Validate() != nil {
		t.Error("note should never fail validation")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestNoteField`
Expected: FAIL - NewNote not defined

**Step 3: Write minimal implementation**

```go
// form/note.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/2389-research/tux/theme"
)

// NoteField is a display-only field for showing information.
type NoteField struct {
	id      string
	title   string
	content string
	focused bool
}

// NewNote creates a new note field.
func NewNote() *NoteField {
	return &NoteField{}
}

// ID sets the field ID.
func (f *NoteField) ID(id string) *NoteField {
	f.id = id
	return f
}

// Title sets the note title.
func (f *NoteField) Title(title string) *NoteField {
	f.title = title
	return f
}

// Content sets the note content.
func (f *NoteField) Content(content string) *NoteField {
	f.content = content
	return f
}

// Field interface implementation

func (f *NoteField) ID_Get() string    { return f.id }
func (f *NoteField) Label_Get() string { return f.title }

func (f *NoteField) Value() any {
	return nil // Notes have no value
}

func (f *NoteField) SetValue(v any) {
	// Notes don't store values
}

func (f *NoteField) Focused() bool {
	return f.focused
}

func (f *NoteField) Focus() {
	f.focused = true
}

func (f *NoteField) Blur() {
	f.focused = false
}

func (f *NoteField) Validate() error {
	return nil // Notes don't validate
}

func (f *NoteField) Init() tea.Cmd {
	return nil
}

func (f *NoteField) HandleKey(key tea.KeyMsg) bool {
	return false // Notes don't handle input
}

func (f *NoteField) Render(width int, th theme.Theme, focused bool) string {
	styles := th.Styles()

	var output string
	if f.title != "" {
		output = styles.Title.Render(f.title) + "\n"
	}
	if f.content != "" {
		output += styles.Body.Render(f.content)
	}

	return lipgloss.NewStyle().Width(width).Render(output)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/note.go form/note_test.go
git commit -m "feat(form): add Note field (display only)"
```

---

### Task 8: Form Compositor

**Files:**
- Modify: `form/form.go`
- Create: `form/form_test.go`

**Step 1: Write the failing test**

```go
// form/form_test.go
package form

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestFormCreation(t *testing.T) {
	f := New(
		NewInput().ID("name").Label("Name"),
		NewSelect[string]().ID("role").Label("Role").Options(
			Option("Admin", "admin"),
			Option("User", "user"),
		),
	)

	if f == nil {
		t.Fatal("form should not be nil")
	}
}

func TestFormValues(t *testing.T) {
	nameField := NewInput().ID("name")
	roleField := NewSelect[string]().ID("role").Options(
		Option("Admin", "admin"),
	).Default("admin")

	f := New(nameField, roleField)
	nameField.SetValue("Alice")

	values := f.Values()
	if values.String("name") != "Alice" {
		t.Errorf("expected name 'Alice', got %s", values.String("name"))
	}
	if values.String("role") != "admin" {
		t.Errorf("expected role 'admin', got %s", values.String("role"))
	}
}

func TestFormNavigation(t *testing.T) {
	f := New(
		NewInput().ID("a"),
		NewInput().ID("b"),
		NewInput().ID("c"),
	)
	f.Init()

	// First field should be focused
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0, got %d", f.FocusedIndex())
	}

	// Tab moves forward
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})
	if f.FocusedIndex() != 1 {
		t.Errorf("expected focused index 1 after Tab, got %d", f.FocusedIndex())
	}

	// Shift+Tab moves back
	f.HandleKey(tea.KeyMsg{Type: tea.KeyShiftTab})
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0 after Shift+Tab, got %d", f.FocusedIndex())
	}
}

func TestFormState(t *testing.T) {
	f := New(NewInput().ID("test"))
	f.Init()

	if f.State() != StateActive {
		t.Error("initial state should be Active")
	}

	// Escape cancels
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})
	if f.State() != StateCancelled {
		t.Error("state should be Cancelled after Escape")
	}
}

func TestFormRender(t *testing.T) {
	f := New(
		NewInput().ID("name").Label("Name"),
	).WithTheme(theme.Get("dracula"))
	f.Init()

	output := f.Render(60, 20)
	if output == "" {
		t.Error("render should produce output")
	}
}

func TestFormGroups(t *testing.T) {
	f := New(
		Group("Page 1",
			NewInput().ID("a"),
		),
		Group("Page 2",
			NewInput().ID("b"),
		),
	)
	f.Init()

	if f.GroupCount() != 2 {
		t.Errorf("expected 2 groups, got %d", f.GroupCount())
	}
	if f.CurrentGroup() != 0 {
		t.Error("should start on group 0")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./form/... -v -run TestForm`
Expected: FAIL - New not defined

**Step 3: Write minimal implementation**

```go
// form/form.go (replace existing content)
package form

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/2389-research/tux/theme"
)

// State represents the form state.
type State int

const (
	StateActive State = iota
	StateSubmitted
	StateCancelled
)

// FormField is the interface that all fields must implement.
type FormField interface {
	ID_Get() string
	Label_Get() string
	Value() any
	SetValue(any)
	Focused() bool
	Focus()
	Blur()
	Validate() error
	Init() tea.Cmd
	HandleKey(key tea.KeyMsg) bool
	Render(width int, th theme.Theme, focused bool) string
}

// FieldGroup represents a group of fields (a page in multi-page forms).
type FieldGroup struct {
	title  string
	fields []FormField
}

// Group creates a field group with a title.
func Group(title string, fields ...FormField) *FieldGroup {
	return &FieldGroup{
		title:  title,
		fields: fields,
	}
}

// Form manages a collection of fields.
type Form struct {
	groups       []*FieldGroup
	currentGroup int
	focusedIndex int
	state        State
	theme        theme.Theme
	onSubmit     func(Values)
	onCancel     func()
}

// New creates a new form. Accepts fields or groups.
func New(items ...any) *Form {
	f := &Form{
		state: StateActive,
		theme: theme.Get("dracula"), // default theme
	}

	// Collect fields into a default group, or use provided groups
	var defaultFields []FormField
	for _, item := range items {
		switch v := item.(type) {
		case *FieldGroup:
			f.groups = append(f.groups, v)
		case FormField:
			defaultFields = append(defaultFields, v)
		}
	}

	// If we have loose fields, put them in a default group
	if len(defaultFields) > 0 {
		f.groups = append([]*FieldGroup{{fields: defaultFields}}, f.groups...)
	}

	// Ensure at least one group
	if len(f.groups) == 0 {
		f.groups = []*FieldGroup{{}}
	}

	return f
}

// WithTheme sets the form theme.
func (f *Form) WithTheme(th theme.Theme) *Form {
	f.theme = th
	return f
}

// OnSubmit sets the submit callback.
func (f *Form) OnSubmit(fn func(Values)) *Form {
	f.onSubmit = fn
	return f
}

// OnCancel sets the cancel callback.
func (f *Form) OnCancel(fn func()) *Form {
	f.onCancel = fn
	return f
}

// Init initializes the form and focuses the first field.
func (f *Form) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, g := range f.groups {
		for _, field := range g.fields {
			if cmd := field.Init(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Focus first field
	if fields := f.currentFields(); len(fields) > 0 {
		fields[0].Focus()
	}

	return tea.Batch(cmds...)
}

// State returns the current form state.
func (f *Form) State() State {
	return f.state
}

// Values returns all field values.
func (f *Form) Values() Values {
	v := make(Values)
	for _, g := range f.groups {
		for _, field := range g.fields {
			if id := field.ID_Get(); id != "" {
				v[id] = field.Value()
			}
		}
	}
	return v
}

// FocusedIndex returns the index of the focused field in the current group.
func (f *Form) FocusedIndex() int {
	return f.focusedIndex
}

// GroupCount returns the number of groups.
func (f *Form) GroupCount() int {
	return len(f.groups)
}

// CurrentGroup returns the current group index.
func (f *Form) CurrentGroup() int {
	return f.currentGroup
}

func (f *Form) currentFields() []FormField {
	if f.currentGroup < len(f.groups) {
		return f.groups[f.currentGroup].fields
	}
	return nil
}

// HandleKey processes keyboard input.
func (f *Form) HandleKey(key tea.KeyMsg) bool {
	if f.state != StateActive {
		return false
	}

	fields := f.currentFields()

	switch key.Type {
	case tea.KeyEscape:
		f.state = StateCancelled
		if f.onCancel != nil {
			f.onCancel()
		}
		return true

	case tea.KeyTab, tea.KeyDown:
		if len(fields) > 0 {
			fields[f.focusedIndex].Blur()
			f.focusedIndex = (f.focusedIndex + 1) % len(fields)
			fields[f.focusedIndex].Focus()
		}
		return true

	case tea.KeyShiftTab, tea.KeyUp:
		if len(fields) > 0 {
			fields[f.focusedIndex].Blur()
			f.focusedIndex--
			if f.focusedIndex < 0 {
				f.focusedIndex = len(fields) - 1
			}
			fields[f.focusedIndex].Focus()
		}
		return true

	case tea.KeyEnter:
		// On last field, submit or go to next group
		if f.focusedIndex == len(fields)-1 {
			if f.currentGroup == len(f.groups)-1 {
				// Last group - submit
				f.state = StateSubmitted
				if f.onSubmit != nil {
					f.onSubmit(f.Values())
				}
			} else {
				// Next group
				f.currentGroup++
				f.focusedIndex = 0
				if newFields := f.currentFields(); len(newFields) > 0 {
					newFields[0].Focus()
				}
			}
			return true
		}
		// Move to next field
		fields[f.focusedIndex].Blur()
		f.focusedIndex++
		fields[f.focusedIndex].Focus()
		return true
	}

	// Delegate to focused field
	if len(fields) > 0 && f.focusedIndex < len(fields) {
		return fields[f.focusedIndex].HandleKey(key)
	}

	return false
}

// Render renders the form.
func (f *Form) Render(width, height int) string {
	if f.theme == nil {
		f.theme = theme.Get("dracula")
	}

	fields := f.currentFields()
	var parts []string

	// Group title
	if f.currentGroup < len(f.groups) && f.groups[f.currentGroup].title != "" {
		title := f.theme.Styles().Title.Render(f.groups[f.currentGroup].title)
		parts = append(parts, title, "")
	}

	// Fields
	for i, field := range fields {
		focused := i == f.focusedIndex
		parts = append(parts, field.Render(width-4, f.theme, focused))
	}

	// Page indicator for multi-group forms
	if len(f.groups) > 1 {
		indicator := f.theme.Styles().Muted.Render(
			strings.Repeat("○ ", f.currentGroup) + "● " + strings.Repeat("○ ", len(f.groups)-f.currentGroup-1),
		)
		parts = append(parts, "", indicator)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./form/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add form/form.go form/form_test.go
git commit -m "feat(form): add Form compositor with groups"
```

---

### Task 9: FormModal Adapter

**Files:**
- Create: `modal/form.go`
- Create: `modal/form_test.go`

**Step 1: Write the failing test**

```go
// modal/form_test.go
package modal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/form"
)

func TestFormModal(t *testing.T) {
	f := form.New(
		form.NewInput().ID("name").Label("Name"),
	)

	m := NewFormModal(FormModalConfig{
		ID:    "test-form",
		Title: "Test Form",
		Form:  f,
	})

	if m.ID() != "test-form" {
		t.Errorf("expected ID 'test-form', got %s", m.ID())
	}
	if m.Title() != "Test Form" {
		t.Errorf("expected title 'Test Form', got %s", m.Title())
	}
}

func TestFormModalSize(t *testing.T) {
	f := form.New(form.NewInput().ID("test"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		Size: SizeLarge,
	})

	if m.Size() != SizeLarge {
		t.Errorf("expected SizeLarge, got %v", m.Size())
	}
}

func TestFormModalOnSubmit(t *testing.T) {
	var submitted bool
	var values form.Values

	f := form.New(form.NewInput().ID("name"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		OnSubmit: func(v form.Values) {
			submitted = true
			values = v
		},
	})

	m.OnPush(80, 24)

	// Simulate Enter to submit (single field form)
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})

	if !submitted {
		t.Error("OnSubmit should have been called")
	}
	if values == nil {
		t.Error("values should not be nil")
	}
}

func TestFormModalOnCancel(t *testing.T) {
	var cancelled bool

	f := form.New(form.NewInput().ID("name"))

	m := NewFormModal(FormModalConfig{
		Form: f,
		OnCancel: func() {
			cancelled = true
		},
	})

	m.OnPush(80, 24)

	// Simulate Escape to cancel
	m.HandleKey(tea.KeyMsg{Type: tea.KeyEscape})

	if !cancelled {
		t.Error("OnCancel should have been called")
	}
}

func TestFormModalRender(t *testing.T) {
	f := form.New(form.NewInput().ID("name").Label("Name"))

	m := NewFormModal(FormModalConfig{
		Title: "Test",
		Form:  f,
	})

	m.OnPush(80, 24)
	output := m.Render(60, 20)

	if output == "" {
		t.Error("render should produce output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./modal/... -v -run TestFormModal`
Expected: FAIL - NewFormModal not defined

**Step 3: Write minimal implementation**

```go
// modal/form.go
package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/2389-research/tux/form"
	"github.com/2389-research/tux/theme"
)

// FormModalConfig configures a FormModal.
type FormModalConfig struct {
	ID       string
	Title    string
	Form     *form.Form
	Size     Size
	Theme    theme.Theme
	OnSubmit func(form.Values)
	OnCancel func()
}

// FormModal adapts a Form to work as a Modal.
type FormModal struct {
	id       string
	title    string
	form     *form.Form
	size     Size
	theme    theme.Theme
	onSubmit func(form.Values)
	onCancel func()
	width    int
	height   int

	boxStyle   lipgloss.Style
	titleStyle lipgloss.Style
}

// NewFormModal creates a new form modal.
func NewFormModal(cfg FormModalConfig) *FormModal {
	size := cfg.Size
	if size == 0 {
		size = SizeMedium
	}

	th := cfg.Theme
	if th == nil {
		th = theme.Get("dracula")
	}

	m := &FormModal{
		id:       cfg.ID,
		title:    cfg.Title,
		form:     cfg.Form,
		size:     size,
		theme:    th,
		onSubmit: cfg.OnSubmit,
		onCancel: cfg.OnCancel,
	}

	// Set up styles
	m.boxStyle = th.Styles().ModalBox
	m.titleStyle = th.Styles().ModalTitle

	// Wire form callbacks
	if m.form != nil {
		m.form.OnSubmit(func(v form.Values) {
			if m.onSubmit != nil {
				m.onSubmit(v)
			}
		})
		m.form.OnCancel(func() {
			if m.onCancel != nil {
				m.onCancel()
			}
		})
		m.form.WithTheme(th)
	}

	return m
}

// ID implements Modal.
func (m *FormModal) ID() string {
	if m.id != "" {
		return m.id
	}
	return "form-modal"
}

// Title implements Modal.
func (m *FormModal) Title() string {
	return m.title
}

// Size implements Modal.
func (m *FormModal) Size() Size {
	return m.size
}

// OnPush implements Modal.
func (m *FormModal) OnPush(width, height int) {
	m.width = width
	m.height = height
	if m.form != nil {
		m.form.Init()
	}
}

// OnPop implements Modal.
func (m *FormModal) OnPop() {}

// HandleKey implements Modal.
func (m *FormModal) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	if m.form == nil {
		return false, nil
	}

	handled := m.form.HandleKey(key)

	// Check if form completed
	switch m.form.State() {
	case form.StateSubmitted:
		return true, func() tea.Msg { return PopMsg{} }
	case form.StateCancelled:
		return true, func() tea.Msg { return PopMsg{} }
	}

	return handled, nil
}

// Render implements Modal.
func (m *FormModal) Render(width, height int) string {
	var content string

	// Title
	if m.title != "" {
		content = m.titleStyle.Render(m.title) + "\n\n"
	}

	// Form content
	if m.form != nil {
		content += m.form.Render(width-6, height-6)
	}

	return m.boxStyle.Width(width - 4).Render(content)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./modal/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add modal/form.go modal/form_test.go
git commit -m "feat(modal): add FormModal adapter"
```

---

### Task 10: Final Integration Test & Coverage

**Files:**
- Create: `form/integration_test.go`

**Step 1: Write integration test**

```go
// form/integration_test.go
package form

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/2389-research/tux/theme"
)

func TestFullFormWorkflow(t *testing.T) {
	var result Values

	f := New(
		NewInput().ID("username").Label("Username").Validate(Required(), MinLength(3)),
		NewSelect[string]().ID("role").Label("Role").Options(
			Option("Admin", "admin"),
			Option("User", "user"),
		).Default("user"),
		NewConfirm().ID("active").Label("Active?"),
	).WithTheme(theme.Get("neo-terminal")).OnSubmit(func(v Values) {
		result = v
	})

	f.Init()

	// Set username
	inputField := f.currentFields()[0].(*InputField)
	inputField.SetValue("alice")

	// Navigate and submit
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})  // to role
	f.HandleKey(tea.KeyMsg{Type: tea.KeyTab})  // to confirm
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter}) // submit

	if f.State() != StateSubmitted {
		t.Errorf("expected StateSubmitted, got %v", f.State())
	}

	if result.String("username") != "alice" {
		t.Errorf("expected username 'alice', got %s", result.String("username"))
	}
	if result.String("role") != "user" {
		t.Errorf("expected role 'user', got %s", result.String("role"))
	}
}

func TestMultiPageForm(t *testing.T) {
	f := New(
		Group("Step 1",
			NewInput().ID("name").Label("Name"),
		),
		Group("Step 2",
			NewInput().ID("email").Label("Email"),
		),
		Group("Confirm",
			NewConfirm().ID("agree").Label("I agree"),
		),
	)

	f.Init()

	if f.GroupCount() != 3 {
		t.Errorf("expected 3 groups, got %d", f.GroupCount())
	}

	// Submit first page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.CurrentGroup() != 1 {
		t.Errorf("expected group 1, got %d", f.CurrentGroup())
	}

	// Submit second page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.CurrentGroup() != 2 {
		t.Errorf("expected group 2, got %d", f.CurrentGroup())
	}

	// Submit final page
	f.HandleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if f.State() != StateSubmitted {
		t.Error("form should be submitted after last page")
	}
}
```

**Step 2: Run all tests with coverage**

Run: `go test ./form/... ./modal/... -cover`
Expected: PASS with 95%+ coverage

**Step 3: Commit**

```bash
git add form/integration_test.go
git commit -m "test(form): add integration tests"
```

---

### Task 11: Update go.mod and Final Verification

**Step 1: Ensure huh dependency is added**

Run: `go mod tidy`

**Step 2: Run full test suite**

Run: `go test ./... -cover`
Expected: All packages PASS with 95%+ coverage

**Step 3: Final commit**

```bash
git add go.mod go.sum
git commit -m "chore: update dependencies for form package"
```

---

## Summary

| Task | Component | Files |
|------|-----------|-------|
| 1 | Field interface & Values | `form/field.go`, `form/form.go` |
| 2 | Validators | `form/validate.go` |
| 3 | Huh theme bridge | `form/huh_theme.go` |
| 4 | Input field | `form/input.go` |
| 5 | Select field | `form/select.go` |
| 6 | Confirm field | `form/confirm.go` |
| 7 | Note field | `form/note.go` |
| 8 | Form compositor | `form/form.go` |
| 9 | FormModal adapter | `modal/form.go` |
| 10 | Integration tests | `form/integration_test.go` |
| 11 | Dependencies | `go.mod`, `go.sum` |

**Not included (YAGNI):** TextArea, MultiSelect, FilePicker - add later when needed.
