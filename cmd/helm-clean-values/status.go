package main

import (
	"errors"
	"fmt"

	"github.com/sarcaustech/helm-clean-values/pkg/logger"
	"github.com/urfave/cli/v2"
)

var statusCmd = &cli.Command{
	Name:   "status",
	Hidden: true,
	Action: func(cCtx *cli.Context) error {

		requirementsMet := true

		logger := logger.Plain{}

		if val := cCtx.String("helm-bin"); val != "" {
			logger.Info(fmt.Sprintf("(required) helm binary path set: %s", val))
		} else {
			logger.Error("(required) helm binary missing")
			requirementsMet = false
		}

		if val := cCtx.Bool("debug"); val {
			logger.Info("(optional) debug logs enabled")
		} else {
			logger.Info("(optional) debug logs disabled")
		}

		if !requirementsMet {
			return errors.New("minimum requirements not met")
		}
		return nil
	},
}
