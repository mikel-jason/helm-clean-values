package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:                 "helm-clean-values",
	Usage:                "identify unused helm values",
	EnableBashCompletion: true,
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
		log.Fatal(err)
	}
}
