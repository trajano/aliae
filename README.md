# aliae 🌱

Cross platform shell management.

> **THIS IS A PERSONAL FORK.**
> **FOR THE OFFICIAL PROJECT, GO TO [aliae.dev](https://aliae.dev) AND [JanDeDobbeleer/aliae](https://github.com/JanDeDobbeleer/aliae).**
> Inno/WinGet is not supported on this fork, only scoop, homebrew and manual.
> The developer of this personal fork uses *git bash*, *bash*, *zsh*, and *fish*
> normally, so other shells may not get the same treatment.

[![License][license]](LISENCE)

![Release Status][release-status]

[![Release][release-badge]][release]

![GitHub Downloads][gh-downloads]

## Documentation

[![Documentation][docs-badge]][docs]

This repo was made with love using GitKraken.

[![GitKraken shield][kraken]][kraken-ref]
<!-- markdownlint-disable first-header-h1 -->

## Fork summary vs upstream main

This fork contains additional changes on top of
[`JanDeDobbeleer/aliae` `main`](https://github.com/JanDeDobbeleer/aliae/tree/main),
focused on Git Bash support and fork-specific distribution.

- Shell behavior and detection:
  - Added Scoop support and shim-aware shell detection behavior.
  - Added `init --tty-only` support and follow-up fixes for piped/non-TTY detection.
  - Added `alias.type` values `python`/`perl` to run inline interpreter source from `PATH`.
  - Added `script.type` values `python`/`perl` to run inline interpreter source from `PATH`.
  - Added `script.state` (`file`, optional `runEvery`, optional `format`) for run-once or interval-controlled scripts.
- Added `aliae get variables` diagnostics to print TTY checks and resolved template variables.
  - Added `aliae state` commands (`list` default, `clear`) to inspect and clear referenced script state files.
  - Improved shell detection for Git Bash and Scoop shim process trees.
- Added shell resolution trace output through `aliae get variables`.
  - Fixed Git Bash PATH delimiter handling on Windows.
  - Added function description support for Fish (`function --description`) and Nushell
    (doc comment) when `description` is set.
- Template features:
  - Added `fileExists` and `dirExists` template helpers.
  - Added `hasCommand` result caching and `hasCommandNoCache` for uncached checks.
  - Added `setArg` helper for shell-agnostic positional argument assignment.
  - Added `progress` template helper for OSC progress output and reset.
  - Added `.ConfigPath` and `.ConfigDir` template variables.
  - Added `.Env` template map for environment variable access like `{{ .Env.DOTFILES }}`.
  - Added top-level `var` entries with precomputed values and `.Var` template map access.
  - Added `.ShellLike` template variable for bash/zsh/fish/tcsh/pwsh/powershell checks.
  - Added `.Hostname` template variable to expose the system hostname.
  - Added `.WSL` template variable to indicate Windows Subsystem for Linux runtime.
- PATH features:
  - Added `ifExists` option to include only existing path entries.
  - Added `cdpath` support with shell-specific rendering and duplicate suppression.
- Configuration features:
  - Added top-level `stat_timeout` to control filesystem stat/existence timeout.
  - Added top-level `cygpath` mode (`internal` default, `external` optional) to choose conversion backend.
  - Added top-level `cache` (enabled by default) to reuse merged+validated config when source files are unchanged.
  - Added top-level `extends` with short/long syntax, cycle detection, and a depth limit.
  - Added conditional `extends[].if` support to skip individual entries when false.
  - Added top-level `progress` with weighted automatic OSC progress across `alias`, `env`,
    `path`, and `script` (`progress.internal` is read from the root config only and ignored
    in included/extended files).
  - Added `env.isPath` to normalize rendered path-valued env vars to OS-native separators.
  - Added `env.ifExists` (path-only) to export env vars only when the rendered directory exists.
  - Added multiline `env` path fallback: picks the first existing entry when `isPath` is true.
  - Added `aliae get config` to print the fully resolved YAML configuration.
  - Added `aliae validate` to validate raw config YAML against the schema.
  - Added strict unknown-property checks in `aliae validate` while keeping runtime parsing permissive.
- Fork packaging and automation:
  - Pointed installer/docs flows to fork-owned release assets where appropriate.
  - Moved Scoop bucket population to
    [`trajano/scoop-bucket`](https://github.com/trajano/scoop-bucket).
  - Moved Homebrew tap population to
    [`trajano/homebrew-tap`](https://github.com/trajano/homebrew-tap).
  - Updated CI/release workflows for fork ownership and publishing boundaries.
- Documentation and contributor process:
  - Added fork-specific guidance, scoped pre-commit checks, and contribution rules.
  - Added GitHub Pages docs deployment workflow for this fork.

## Join the community

[![Discord][discord]][discord-link]

## ❤ Support ❤

[![GitHub][github-badge]][github-sponsors] - One time support, or a recurring donation?

[![Ko-Fi][kofi-badge]][kofi] - No coffee, no code.

[release-status]: https://img.shields.io/github/actions/workflow/status/jandedobbeleer/aliae/release.yml?branch=main
[release-badge]: https://img.shields.io/github/v/release/jandedobbeleer/aliae?label=Release
[release]: https://github.com/JanDeDobbeleer/aliae/releases/latest
[license]: https://img.shields.io/github/license/JanDeDobbeleer/aliae.svg
[gh-downloads]: https://img.shields.io/github/downloads/jandedobbeleer/aliae/total?color=pink&label=GitHub%20Downloads
[docs-badge]: https://img.shields.io/badge/Docs-GitHub%20Pages-blue
[docs]: https://trajano.github.io/aliae/
[kraken]: https://img.shields.io/badge/GitKraken-Legendary%20Git%20Tools-teal?style=plastic&logo=gitkraken
[kraken-ref]: https://www.gitkraken.com/invite/nQmDPR9D
[discord]: https://img.shields.io/discord/1023597603331526656
[discord-link]: https://discord.gg/n7E3DkXssv
[github-badge]: https://img.shields.io/badge/-Sponsor-fafbfc?logo=GitHub%20Sponsors
[github-sponsors]: https://github.com/sponsors/JanDeDobbeleer
[kofi-badge]: https://img.shields.io/badge/Ko--fi-Buy%20me%20a%20coffee!-%2346b798.svg
[kofi]: https://ko-fi.com/jandedobbeleer
