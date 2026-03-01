package config

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/shell"
)

type schemaBuilder struct {
	aliases []map[string]any
	vars    []map[string]any
	envs    []map[string]any
	paths   []map[string]any
	cdpaths []map[string]any
	scripts []map[string]any
	links   []map[string]any
}

func buildSchemaDocument(aliae *Aliae) map[string]any {
	document := make(map[string]any)
	if aliae == nil {
		return document
	}

	builder := &schemaBuilder{
		aliases: make([]map[string]any, 0, len(aliae.Aliae)),
		vars:    make([]map[string]any, 0, len(aliae.Vars)),
		envs:    make([]map[string]any, 0, len(aliae.Envs)),
		paths:   make([]map[string]any, 0, len(aliae.Paths)),
		cdpaths: make([]map[string]any, 0, len(aliae.CDPaths)),
		scripts: make([]map[string]any, 0, len(aliae.Scripts)),
		links:   make([]map[string]any, 0, len(aliae.Links)),
	}

	WalkConfig(aliae, ConfigVisitorFuncs{
		OnAlias: func(alias *shell.Alias) {
			builder.aliases = append(builder.aliases, aliasSchemaItem(alias))
		},
		OnVar: func(variable *Var) {
			builder.vars = append(builder.vars, varSchemaItem(variable))
		},
		OnEnv: func(env *shell.Env) {
			builder.envs = append(builder.envs, envSchemaItem(env))
		},
		OnPath: func(path *shell.Path) {
			builder.paths = append(builder.paths, pathSchemaItem(path))
		},
		OnCDPath: func(path *shell.CDPath) {
			builder.cdpaths = append(builder.cdpaths, cdpathSchemaItem(path))
		},
		OnScript: func(script *shell.Script) {
			builder.scripts = append(builder.scripts, scriptSchemaItem(script))
		},
		OnLink: func(link *shell.Link) {
			builder.links = append(builder.links, linkSchemaItem(link))
		},
	})

	if len(builder.aliases) > 0 {
		document["alias"] = builder.aliases
	}
	if len(builder.vars) > 0 {
		document["var"] = builder.vars
	}
	if len(builder.envs) > 0 {
		document["env"] = builder.envs
	}
	if len(builder.paths) > 0 {
		document["path"] = builder.paths
	}
	if len(builder.cdpaths) > 0 {
		document["cdpath"] = builder.cdpaths
	}
	if len(builder.scripts) > 0 {
		document["script"] = builder.scripts
	}
	if len(builder.links) > 0 {
		document["link"] = builder.links
	}

	if progress, ok := progressSchemaItem(aliae); ok {
		document["progress"] = progress
	}

	if aliae.StatTimeout > 0 {
		document["stat_timeout"] = aliae.StatTimeout.String()
	}

	if len(strings.TrimSpace(aliae.Cygpath)) > 0 {
		document["cygpath"] = aliae.Cygpath
	}

	if aliae.Cache {
		document["cache"] = true
	}

	return document
}

func varSchemaItem(variable *Var) map[string]any {
	item := map[string]any{
		"name":  variable.Name,
		"value": string(variable.Value),
	}
	if len(variable.If) > 0 {
		item["if"] = string(variable.If)
	}
	return item
}

func aliasSchemaItem(alias *shell.Alias) map[string]any {
	item := map[string]any{
		"name":  alias.Name,
		"value": string(alias.Value),
	}
	if len(alias.Type) > 0 {
		item["type"] = string(alias.Type)
	}
	if len(alias.If) > 0 {
		item["if"] = string(alias.If)
	}
	return item
}

func envSchemaItem(env *shell.Env) map[string]any {
	item := map[string]any{
		"name":  env.Name,
		"value": env.Value,
	}
	if len(env.Delimiter) > 0 {
		item["delimiter"] = string(env.Delimiter)
	}
	if len(env.If) > 0 {
		item["if"] = string(env.If)
	}
	if len(env.Type) > 0 {
		item["type"] = string(env.Type)
	}
	if env.IsPath {
		item["isPath"] = true
	}
	if env.IfExists {
		item["ifExists"] = true
	}
	if env.Persist {
		item["persist"] = true
	}
	return item
}

func pathSchemaItem(path *shell.Path) map[string]any {
	item := map[string]any{
		"value": string(path.Value),
	}
	if len(path.If) > 0 {
		item["if"] = string(path.If)
	}
	if path.Persist {
		item["persist"] = true
	}
	if path.Force {
		item["force"] = true
	}
	if path.IfExists {
		item["ifExists"] = true
	}
	return item
}

func cdpathSchemaItem(path *shell.CDPath) map[string]any {
	item := map[string]any{
		"value": string(path.Value),
	}
	if len(path.If) > 0 {
		item["if"] = string(path.If)
	}
	if path.Force {
		item["force"] = true
	}
	if path.IfExists {
		item["ifExists"] = true
	}
	return item
}

func scriptSchemaItem(script *shell.Script) map[string]any {
	item := map[string]any{
		"value": string(script.Value),
	}
	if len(script.Type) > 0 {
		item["type"] = string(script.Type)
	}
	if len(script.If) > 0 {
		item["if"] = string(script.If)
	}
	if script.Weight > 0 {
		item["weight"] = script.Weight
	}
	if len(script.State.File) > 0 {
		state := map[string]any{
			"file": string(script.State.File),
		}
		if len(script.State.RunEvery) > 0 {
			state["runEvery"] = script.State.RunEvery
		}
		if len(script.State.Format) > 0 {
			state["format"] = string(script.State.Format)
		}
		item["state"] = state
	}
	return item
}

func linkSchemaItem(link *shell.Link) map[string]any {
	item := map[string]any{
		"name":   string(link.Name),
		"target": string(link.Target),
	}
	if len(link.If) > 0 {
		item["if"] = string(link.If)
	}
	if link.MkDir {
		item["mkdir"] = true
	}
	return item
}
