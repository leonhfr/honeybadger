package engine

// OptionType represents an option's type.
type OptionType int

const (
	OptionBoolean OptionType = iota // OptionBoolean represents a boolean option.
	OptionInteger                   // OptionInteger represents an integer option.
	OptionEnum                      // OptionEnum represents an enum option.
)

// Option represents an available option.
type Option struct {
	Type    OptionType
	Name    string
	Default string
	Min     string
	Max     string
	Vars    []string
}

// Options returns the list of available options.
func Options() []Option {
	return nil
}
