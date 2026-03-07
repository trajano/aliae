package appcmd

import (
	"fmt"
	"io"
	"os"

	initpkg "github.com/jandedobbeleer/aliae/src/init"
)

type InitCommand struct {
	Out        io.Writer
	Err        io.Writer
	ConfigPath string
	Shell      string
	Print      bool
	TTYOnly    bool
	StdinTTY   bool
}

func (c InitCommand) Execute() error {
	if c.TTYOnly && !c.StdinTTY {
		return nil
	}

	errWriter := c.Err
	if errWriter == nil {
		errWriter = os.Stderr
	}
	initpkg.SetProgressWriter(errWriter)
	defer initpkg.SetProgressWriter(os.Stderr)

	outWriter := c.Out
	if outWriter == nil {
		outWriter = io.Discard
	}

	script := initpkg.Init(c.ConfigPath, c.Shell, c.Print)
	_, err := fmt.Fprint(outWriter, script)
	return err
}
