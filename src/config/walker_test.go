package config

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestWalkConfigTraversesStableOrder(t *testing.T) {
	cfg := &Aliae{
		Extends: []Extend{
			{Path: "./base.yaml", If: "true"},
		},
		Vars: Vars{
			{Name: "var1"},
			nil,
		},
		Envs: shell.Envs{
			{Name: "ENV1"},
			nil,
		},
		Paths: shell.Paths{
			{Value: "/path1"},
			nil,
		},
		CDPaths: shell.CDPaths{
			{Value: "/cd1"},
			nil,
		},
		Aliae: shell.Aliae{
			{Name: "alias1"},
			nil,
		},
		Links: shell.Links{
			{Name: "link1"},
			nil,
		},
		Scripts: shell.Scripts{
			{Value: "echo script1"},
			nil,
		},
	}

	events := make([]string, 0)
	visitor := ConfigVisitorFuncs{
		OnExtend: func(item *Extend) {
			events = append(events, "extend:"+item.Path)
		},
		OnVar: func(item *Var) {
			events = append(events, "var:"+item.Name)
		},
		OnEnv: func(item *shell.Env) {
			events = append(events, "env:"+item.Name)
		},
		OnPath: func(item *shell.Path) {
			events = append(events, "path:"+string(item.Value))
		},
		OnCDPath: func(item *shell.CDPath) {
			events = append(events, "cdpath:"+string(item.Value))
		},
		OnAlias: func(item *shell.Alias) {
			events = append(events, "alias:"+item.Name)
		},
		OnLink: func(item *shell.Link) {
			events = append(events, "link:"+string(item.Name))
		},
		OnScript: func(item *shell.Script) {
			events = append(events, "script:"+string(item.Value))
		},
	}

	WalkConfig(cfg, visitor)

	assert.Equal(t, []string{
		"extend:./base.yaml",
		"var:var1",
		"env:ENV1",
		"path:/path1",
		"cdpath:/cd1",
		"alias:alias1",
		"link:link1",
		"script:echo script1",
	}, events)
}

func TestWalkConfigNoopForNilInput(t *testing.T) {
	cfg := &Aliae{
		Aliae: shell.Aliae{
			{Name: "alias1"},
		},
	}

	// no panic
	WalkConfig(nil, ConfigVisitorFuncs{})
	WalkConfig(cfg, nil)
}
