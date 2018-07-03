package registry

import (
	"github.com/suzuki-shunsuke/gria/domain"
	"github.com/suzuki-shunsuke/gria/infra"
)

// NewFileWriter returns a new domain.FileWriter .
func NewFileWriter() domain.FileWriter {
	return infra.FileWriter{}
}
