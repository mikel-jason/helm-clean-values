package utils

import (
	"io"
	"os"
)

func ReaderToString(in io.Reader) string {
	stdinBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "" // TODO don't suppress error
	}

	return string(stdinBytes)
}
