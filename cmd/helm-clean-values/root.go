package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var version, commit, date string // filled by goreleaser

var app = &cli.App{
	Name:                 "helm-clean-values",
	Usage:                "identify unused helm values",
	EnableBashCompletion: true,
	Version:              fmt.Sprintf("%s (%s), %s", version, commit, date),
	Flags: []cli.Flag{
		// https://helm.sh/docs/topics/plugins/#environment-variables
		&cli.StringFlag{
			Name:    "helm-bin",
			Usage:   "path to the helm binary",
			EnvVars: []string{"HELM_BIN"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			EnvVars: []string{"HELM_DEBUG"},
		},
	},
	Commands: []*cli.Command{
		mutateCmd,
		simpleCmd,
		statusCmd,
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
