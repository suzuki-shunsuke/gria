package domain

type (
	FileWriter interface {
		Append(dest string, data []byte) error
		Create(dest string, data []byte) error
	}
)
