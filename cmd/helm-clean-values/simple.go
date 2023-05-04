package main

import (
	"fmt"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/providers/helm"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/providers/stdin"
	"github.com/sarcaustech/helm-clean-values/pkg/core/adapters/selectors/simple"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var simpleCmd = &cli.Command{
	Name:  "simple",
	Usage: "detect unused helm values by comparing with the chart's default values",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "stdin",
			Usage: "read input values from STDIN",
		},
		&cli.StringFlag{
			Name:  "chart",
			Usage: "helm prompt to get the chart",
		},
	},
	Before: func(cCtx *cli.Context) error {
		if !cCtx.Bool("stdin") {
			return fmt.Errorf("no input values provided")
		}
		return nil
	},
	Action: func(cCtx *cli.Context) (err error) {

		inputProvider := stdin.ValuesProvider{}
		referenceProvider := helm.ValuesProvider{
			BinaryPath: cCtx.String("helm-bin"),
			Prompt:     cCtx.String("chart"),
		}
		selector := simple.Selector{}

		cleanedValues, err := core.Run(&inputProvider, &referenceProvider, &selector)
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
