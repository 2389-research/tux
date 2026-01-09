# Forms Package Design

> Abstract form system for tux with language-agnostic spec and huh-backed Go implementation.

## Overview

Forms provide rich input collection (text, selections, confirmations) that integrate with tux's modal system. The design is language-agnostic - Go uses huh internally, other languages will use native form libraries.

## Architecture

```
┌─────────────────────────────────────────┐
│ Modal Stack (existing)                   │
│   └── Modal interface                    │
│         └── FormModal (new adapter)      │
│               └── Form                   │
│                     └── Field[]          │
│                           └── Input      │
│                           └── Select     │
│                           └── Confirm    │
│                           └── ...        │
└─────────────────────────────────────────┘
```

**Package structure:**
```
tux/
├── form/           # NEW
│   ├── form.go     # Form struct, builder
│   ├── field.go    # Field interface
│   ├── input.go    # Input field
│   ├── textarea.go # TextArea field
│   ├── select.go   # Select field
│   ├── multiselect.go
│   ├── confirm.go
│   ├── note.go
│   ├── filepicker.go
│   ├── validate.go # Validators
│   └── huh_theme.go # Theme bridge (Go-specific)
├── modal/
│   └── form.go     # FormModal adapter (new file)
```

## Field Types

Seven field types with consistent builder API:

### Input (single-line text)
```go
form.NewInput().
    Label("Username").
    Placeholder("Enter username...").
    Validate(form.Required(), form.MinLength(3))
```

### TextArea (multi-line text)
```go
form.NewTextArea().
    Label("Description").
    Lines(5).
    MaxLength(500)
```

### Select (pick one)
```go
form.NewSelect[string]().
    Label("Theme").
    Options(
        form.Option("Dracula", "dracula"),
        form.Option("Nord", "nord"),
    ).
    Default("dracula")
```

### MultiSelect (pick many)
```go
form.NewMultiSelect[string]().
    Label("Features").
    Options(...).
    Min(1).Max(3)
```

### Confirm (yes/no)
```go
form.NewConfirm().
    Label("Delete this file?").
    Affirmative("Yes, delete").
    Negative("Cancel")
```

### Note (display only)
```go
form.NewNote().
    Title("Warning").
    Content("This action cannot be undone.")
```

### FilePicker (file selection)
```go
form.NewFilePicker().
    Label("Select file").
    AllowedTypes(".go", ".md").
    ShowHidden(false)
```

## Field Interface

```go
type Field interface {
    // Identity
    ID() string
    Label() string

    // State
    Value() any
    SetValue(any)
    Focused() bool

    // Validation
    Validate() error

    // Rendering (theme-aware)
    Render(width int, theme theme.Theme, focused bool) string

    // Input handling
    HandleKey(key tea.KeyMsg) (handled bool)
}
```

## Form Composition

### Simple form (single group)
```go
f := form.New(
    form.NewInput().Label("Name").ID("name"),
    form.NewSelect[string]().Label("Role").ID("role").Options(...),
    form.NewConfirm().Label("Active?").ID("active"),
)
```

### Multi-page form (wizard-style)
```go
f := form.New(
    form.Group("Account",
        form.NewInput().Label("Username").ID("username"),
        form.NewInput().Label("Email").ID("email"),
    ),
    form.Group("Preferences",
        form.NewSelect[string]().Label("Theme").ID("theme").Options(...),
    ),
    form.Group("Confirm",
        form.NewNote().Content("Review your choices."),
        form.NewConfirm().Label("Create account?").ID("confirm"),
    ),
)
```

### Form struct
```go
type Form struct {
    groups       []*Group
    currentGroup int
    theme        theme.Theme
    onSubmit     func(Values)
    onCancel     func()
}

type Group struct {
    title  string
    fields []Field
}

type Values map[string]any

func (v Values) String(id string) string
func (v Values) Bool(id string) bool
func (v Values) Strings(id string) []string
```

### Navigation
- Tab/Shift+Tab or ↓/↑: move between fields
- Enter on last field: next group (or submit if last group)
- Esc: cancel form
- Page indicators for multi-group forms

## Modal Integration

FormModal adapts Form to work in the modal stack:

```go
type FormModal struct {
    id       string
    title    string
    form     *form.Form
    size     Size
    onSubmit func(form.Values)
    onCancel func()
}

// Usage
modal := modal.NewFormModal(modal.FormModalConfig{
    ID:       "settings-form",
    Title:    "Settings",
    Form:     myForm,
    Size:     modal.SizeMedium,
    OnSubmit: func(v form.Values) {
        theme := v.String("theme")
    },
})
manager.Push(modal)
```

FormModal implements the Modal interface, delegating to the form for key handling and rendering, and auto-popping on submit/cancel.

## Theming

Forms use tux themes directly:

```go
f := form.New(...).WithTheme(theme.Get("dracula"))
```

Fields render using theme styles:
- `Input` / `InputFocused` - text input boxes
- `ListItem` / `ListItemSelected` - select options
- `Button` / `ButtonActive` - confirm buttons
- `Error` - validation errors
- `Muted` - placeholders, hints

No changes needed to existing theme package.

## Validation

Built-in validators with composable API:

```go
// Built-in
form.Required()
form.MinLength(n)
form.MaxLength(n)
form.Pattern(regex, msg)
form.Email()
form.MinSelected(n)
form.MaxSelected(n)

// Custom
form.NewInput().Validate(func(v any) error {
    if strings.Contains(v.(string), " ") {
        return errors.New("no spaces allowed")
    }
    return nil
})
```

**Timing:**
- On field blur (moving to next field)
- On submit attempt
- Errors display inline below field
- Form won't submit until all fields valid

## Go Implementation

Uses huh internally - apps never import huh directly:

```go
type InputField struct {
    id          string
    label       string
    value       string
    validators  []Validator
    huhField    *huh.Input  // internal
}
```

Theme bridge converts tux theme to huh theme:

```go
func toHuhTheme(th theme.Theme) *huh.Theme {
    t := huh.ThemeBase()
    t.Focused.Title = lipgloss.NewStyle().
        Foreground(th.Primary()).Bold(true)
    // ... map all styles
    return t
}
```

## Testing Strategy

1. **Unit tests per field** - value get/set, rendering
2. **Validation tests** - each validator, composition
3. **Form tests** - navigation, grouping, state
4. **Modal integration** - submit/cancel flow

Coverage target: 95%+

## Multi-Language Note

This design is language-agnostic:
- **Go:** uses huh internally
- **Rust:** will use ratatui widgets
- **TypeScript:** will use ink form components
- **Python:** will use Textual widgets

The abstract Field/Form API remains consistent across implementations.
