// Package bounds provides types to represent a bounds check
package bounds

// Result enumerates the result of a bounds check
type Result int

// Result values
const (
	LessThan Result = iota
	Inside
	GreaterThan
)
