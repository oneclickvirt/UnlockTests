package utils

import (
	"io"

	"github.com/mattn/go-colorable"
)

var (
	ColorStdout io.Writer = colorable.NewColorableStdout()
	ColorStderr io.Writer = colorable.NewColorableStderr()
)
