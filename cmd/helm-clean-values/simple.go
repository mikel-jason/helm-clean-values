package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sarcaustech/helm-clean-values/pkg/core"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var simpleCmd = &cli.Command{
	Name:  "simple",
	Usage: "detect unused helm values by comparing with the chart's default values",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "chart",
			Usage: "helm prompt to get the chart",
		},
	},
	Action: func(cCtx *cli.Context) (err error) {

		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("cannot read stdin: %w", err)
		}
		var values map[string]interface{}
		if err = yaml.Unmarshal(stdinBytes, &values); err != nil {
			return fmt.Errorf("cannot parse stdin to YAML: %w", err)
		}

		chartPrompt := cCtx.String("chart")

		cmd := exec.Command(cCtx.String("helm-bin"), "show", "values", chartPrompt)
		defaultValuesBytes, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("helm show values %s failed: %w", chartPrompt, err)
		}

		var defaultValues map[string]interface{}
		if err = yaml.Unmarshal(defaultValuesBytes, &defaultValues); err != nil {
			return fmt.Errorf("cannot parse stdin to YAML: %w", err)
		}

		masked := core.Mask(values, defaultValues)
		bytes, err := yaml.Marshal(masked)
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))

		return nil
	},
}
