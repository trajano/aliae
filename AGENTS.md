# Git Bash support

This fork exists to keep `aliae` usable for Git Bash users.

## Fork publishing

Do not publish to the upstream repository.
Use `trajano/aliae` as the publishing target for release assets and workflow automation.

## Commitlint requirement

Every commit message must pass commitlint.
Use Conventional Commit format and allowed types from `.commitlintrc.yml`:

- `chore`
- `feat`
- `fix`
- `docs`
- `refactor`
- `perf`
- `test`
- `revert`
- `ci`

## Rebase requirement

Before treating a branch as a merge candidate, the agent must rebase it onto `origin/HEAD`.

## Pre-commit checks

Run checks based on the files changed in the commit.

If any file under `src/**` changed, run these commands from the repo root:

`cd src && go fmt ./...`

`cd src && go mod tidy`

`cd src && go test ./...`

If any file under `website/**` changed, run this command from the repo root:

`cd website && npm ci && npm run build`

If any file under `bucket/**` or `packages/scoop/**` changed, run this command from the repo root:

`Get-Content bucket/aliae.template.json | ConvertFrom-Json | Out-Null`

`Get-Content bucket/aliae.json | ConvertFrom-Json | Out-Null`

If multiple scopes changed, run all applicable check sets.

If neither `src/**` nor `website/**` changed, these scoped pre-commit checks are not required.

## Website documentation requirement

When a new feature is added, the agent must also update the website documentation in the same change.

## Template variable change checklist

When adding or changing template variables, update all of the following in the same change:

- Runtime/template variable source in `src/context/runtime.go` and related config wiring in `src/config/config.go`.
- Template coverage tests in `src/config/init_test.go` and helper/unit tests in `src/config/config_test.go` or `src/shell/template_test.go`.
- Website documentation in `website/docs/setup/templates.mdx`.

## PR update requirement

When updating an open PR, the agent must ensure the PR description still matches the current change scope.
The agent must also add an update comment on the PR summarizing what changed in the update.
