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
  - Added `aliae debug` diagnostics to print TTY checks and resolved template variables.
  - Improved shell detection for Git Bash and Scoop shim process trees.
  - Added shell resolution trace output through `aliae debug`.
  - Fixed Git Bash PATH delimiter handling on Windows.
  - Added function description support for Fish (`function --description`) and Nushell
    (doc comment) when `description` is set.
- Template features:
  - Added `fileExists` and `dirExists` template helpers.
  - Added `progress` template helper for OSC progress output and reset.
  - Added `.ConfigPath` and `.ConfigDir` template variables.
  - Added `.Hostname` template variable to expose the system hostname.
  - Added `.WSL` template variable to indicate Windows Subsystem for Linux runtime.
- PATH features:
  - Added `ifExists` option to include only existing path entries.
  - Added `cdpath` support with shell-specific rendering and duplicate suppression.
- Configuration features:
  - Added top-level `stat_timeout` to control filesystem stat/existence timeout.
- Fork packaging and automation:
  - Pointed installer/docs flows to fork-owned release assets where appropriate.
  - Moved Scoop/Homebrew bucket and tap population to dedicated external repositories.
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
