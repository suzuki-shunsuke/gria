package registry

import (
	"github.com/suzuki-shunsuke/gria/domain"
	"github.com/suzuki-shunsuke/gria/infra"
)

func NewFileWriter() domain.FileWriter {
	return infra.FileWriter{}
}
