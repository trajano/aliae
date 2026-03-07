package appcmd

import cfg "github.com/jandedobbeleer/aliae/src/config"

type ValidateCommand struct {
	ConfigPath string
}

func (c ValidateCommand) Execute() error {
	return cfg.ValidateConfig(c.ConfigPath)
}
