package logger

import (
	"fmt"
	"os"
)

type Plain struct {
	EnableDebug bool
}

func (p *Plain) Error(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR\t%s\n", msg)
}

func (p *Plain) Warn(msg string) {
	fmt.Fprintf(os.Stderr, "WARN\t%s\n", msg)
}

func (p *Plain) Info(msg string) {
	fmt.Fprintf(os.Stderr, "INFO\t%s\n", msg)
}

func (p *Plain) Debug(msg string) {
	if !p.EnableDebug {
		return
	}
	fmt.Fprintf(os.Stderr, "DEBUG\t%s\n", msg)
}
