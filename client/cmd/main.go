package main

import (
	"github.com/ajugalushkin/goph-keeper/internal/cli"
)

func main() {
	cli.Execute(cli.NewRootCmd())
}
