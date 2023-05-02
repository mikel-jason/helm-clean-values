package main

import (
	"errors"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var statusCmd = &cli.Command{
	Name:   "status",
	Hidden: true,
	Action: func(cCtx *cli.Context) error {

		var required = slog.Bool("required", true)
		var optional = slog.Bool("required", false)

		requirementsMet := true

		if val := cCtx.String("helm-bin"); val != "" {
			slog.Info(
				"helm binary path set",
				slog.String("value", val),
				required,
			)
		} else {
			slog.Error(
				"helm binary not path set",
				slog.String("value", val),
				required,
			)
			requirementsMet = false
		}

		if val := cCtx.Bool("debug"); val {
			slog.Info(
				"debug logs enabled",
				optional,
			)
		} else {
			slog.Info(
				"debug logs disabled",
				optional,
			)
		}

		if !requirementsMet {
			return errors.New("minimum requirements not met")
		}
		return nil
	},
}
