package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/gria/command"
	"github.com/suzuki-shunsuke/gria/domain"
)

func main() {
	app := cli.NewApp()
	app.Name = "gria"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "generate test function's scaffold for golang"
	app.Action = command.Gen
	app.Run(os.Args)
}
