package main

import (
	"fmt"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/providers/file"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/providers/helm"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/providers/stdin"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/selectors/simple"
	"github.com/sarcaustech/helm-clean-values/pkg/logger"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var simpleInputMethod string // centralize once more than one selector cmd

var simpleCmd = &cli.Command{
	Name:  "simple",
	Usage: "detect unused helm values by comparing with the chart's default values",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "stdin",
			Usage: "read input values from STDIN",
		},
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "filepath to read input values from",
		},
		&cli.StringFlag{
			Name:  "chart",
			Usage: "helm prompt to get the chart",
		},
	},
	Before: func(cCtx *cli.Context) error {
		methods := 0

		if cCtx.Bool("stdin") {
			methods++
			simpleInputMethod = "stdin"
		}
		if cCtx.String("file") != "" {
			methods++
			simpleInputMethod = "file"
		}

		if methods > 1 {
			return fmt.Errorf("too input values provided, expected only one method")
		}
		if methods == 0 {
			return fmt.Errorf("no input values provided")
		}

		return nil
	},
	Action: func(cCtx *cli.Context) (err error) {

		logger := &logger.Plain{
			EnableDebug: cCtx.Bool("debug"),
		}

		var inputProvider core.ValuesProvider
		switch simpleInputMethod {
		case "stdin":
			inputProvider = &stdin.ValuesProvider{}
		case "file":
			inputProvider = &file.ValuesProvider{
				Path: cCtx.String("file"),
			}
		default:
			return fmt.Errorf("unknown input method %s", simpleInputMethod)
		}

		referenceProvider := &helm.ValuesProvider{
			BinaryPath: cCtx.String("helm-bin"),
			Prompt:     cCtx.String("chart"),
		}
		selector := &simple.Selector{}

		cleanedValues, err := core.Run(logger, inputProvider, referenceProvider, selector)
		if err != nil {
			return err
		}
		bytes, err := yaml.Marshal(cleanedValues)
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))

		return nil
	},
}
