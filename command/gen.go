package command

import (
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/gria/domain"
	"github.com/suzuki-shunsuke/gria/usecase"
)

// Gen is the command entrypoint.
func Gen(c *cli.Context) error {
	args := domain.GenArgs{
		Targets: []string(c.Args()),
	}
	err := usecase.Gen(args)
	if err == nil {
		return nil
	}
	return cli.NewExitError(err.Error(), 1)
}
