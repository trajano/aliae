package init

import (
	"io"

	cfg "github.com/jandedobbeleer/aliae/src/config"
)

func SetProgressWriter(writer io.Writer) {
	cfg.SetInitProgressWriter(writer)
}
