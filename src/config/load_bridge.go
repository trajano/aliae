package config

import cfgload "github.com/jandedobbeleer/aliae/src/config/load"

func setTemplateConfigContext(configPath string) {
	cfgload.SetTemplateConfigContext(configPath)
}

func resolveConfigDir(configPath string) string {
	return cfgload.ResolveConfigDir(configPath)
}

func home() string {
	return cfgload.Home()
}

func resolveConfigPath(configPath string) string {
	return cfgload.ResolveConfigPath(configPath)
}

func ResolveTemplateContext(configPath string) (string, string) {
	return cfgload.ResolveTemplateContext(configPath)
}
