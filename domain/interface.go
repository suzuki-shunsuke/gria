package domain

type (
	// FileWriter represents the interface to add tests to codes.
	FileWriter interface {
		Append(dest string, data []byte) error
		Create(dest string, data []byte) error
	}
)
