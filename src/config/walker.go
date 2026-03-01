package config

import "github.com/jandedobbeleer/aliae/src/shell"

// ConfigVisitor receives configuration items in a stable traversal order.
type ConfigVisitor interface {
	VisitExtend(*Extend)
	VisitVar(*Var)
	VisitEnv(*shell.Env)
	VisitPath(*shell.Path)
	VisitCDPath(*shell.CDPath)
	VisitAlias(*shell.Alias)
	VisitLink(*shell.Link)
	VisitScript(*shell.Script)
}

// ConfigVisitorFuncs is an adapter for ConfigVisitor where each callback is optional.
type ConfigVisitorFuncs struct {
	OnExtend func(*Extend)
	OnVar    func(*Var)
	OnEnv    func(*shell.Env)
	OnPath   func(*shell.Path)
	OnCDPath func(*shell.CDPath)
	OnAlias  func(*shell.Alias)
	OnLink   func(*shell.Link)
	OnScript func(*shell.Script)
}

func (v ConfigVisitorFuncs) VisitExtend(item *Extend) {
	if v.OnExtend != nil {
		v.OnExtend(item)
	}
}

func (v ConfigVisitorFuncs) VisitVar(item *Var) {
	if v.OnVar != nil {
		v.OnVar(item)
	}
}

func (v ConfigVisitorFuncs) VisitEnv(item *shell.Env) {
	if v.OnEnv != nil {
		v.OnEnv(item)
	}
}

func (v ConfigVisitorFuncs) VisitPath(item *shell.Path) {
	if v.OnPath != nil {
		v.OnPath(item)
	}
}

func (v ConfigVisitorFuncs) VisitCDPath(item *shell.CDPath) {
	if v.OnCDPath != nil {
		v.OnCDPath(item)
	}
}

func (v ConfigVisitorFuncs) VisitAlias(item *shell.Alias) {
	if v.OnAlias != nil {
		v.OnAlias(item)
	}
}

func (v ConfigVisitorFuncs) VisitLink(item *shell.Link) {
	if v.OnLink != nil {
		v.OnLink(item)
	}
}

func (v ConfigVisitorFuncs) VisitScript(item *shell.Script) {
	if v.OnScript != nil {
		v.OnScript(item)
	}
}

// WalkConfig traverses a fully loaded configuration in a stable section order.
func WalkConfig(cfg *Aliae, visitor ConfigVisitor) {
	if cfg == nil || visitor == nil {
		return
	}

	for i := range cfg.Extends {
		visitor.VisitExtend(&cfg.Extends[i])
	}

	for _, item := range cfg.Vars {
		if item != nil {
			visitor.VisitVar(item)
		}
	}

	for _, item := range cfg.Envs {
		if item != nil {
			visitor.VisitEnv(item)
		}
	}

	for _, item := range cfg.Paths {
		if item != nil {
			visitor.VisitPath(item)
		}
	}

	for _, item := range cfg.CDPaths {
		if item != nil {
			visitor.VisitCDPath(item)
		}
	}

	for _, item := range cfg.Aliae {
		if item != nil {
			visitor.VisitAlias(item)
		}
	}

	for _, item := range cfg.Links {
		if item != nil {
			visitor.VisitLink(item)
		}
	}

	for _, item := range cfg.Scripts {
		if item != nil {
			visitor.VisitScript(item)
		}
	}
}
