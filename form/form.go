// form/form.go
package form

// State represents the form state.
type State int

const (
	StateActive State = iota
	StateSubmitted
	StateCancelled
)
