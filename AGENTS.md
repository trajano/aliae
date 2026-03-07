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

Use grouped validation by changed scope:

- `src/**` changes:
  - `cd src && go fmt ./...`
  - `cd src && go mod tidy`
  - `cd src && go test ./...`
  - `cd src && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run`
  - `cd src && go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest`
  - `cd src && "$(go env GOPATH)"/bin/fieldalignment ./...`
- `website/**` changes:
  - `cd website && npm ci && npm run build`
- Markdown/MDX changes (`*.md`, `*.mdx`):
  - `npx --yes markdownlint-cli $(rg --files -g "*.md" -g "*.mdx")`

If multiple scopes changed, run all applicable groups.

If neither `src/**` nor `website/**` changed, those scoped groups are not required.

Field alignment is part of definition of done for `src/**` changes: local `fieldalignment ./...` must pass before pushing.

## Website documentation requirement

When a new feature is added, the agent must also update the website documentation in the same change.
For features added in this fork since `https://github.com/JanDeDobbeleer/aliae`
`main`, update both `README.md` and `website/docs/introduction.mdx` as part of PR definition of done.

## Testing

- When testing `init` from the agent, `--tty-only=false` must be set the TTY checks will fail when running inside the agent.

## Architecture review notes

- Script `if` eligibility is intentionally frozen during init progress-weight computation and reused during script render in the same run. This keeps progress accounting and applied elements consistent. Do not flag this as a bug in architecture reviews unless the behavior explicitly changes.
- aliae CLI execution is intentionally short-lived with a startup target under one second. Process-global runtime/output/cache state in `src/**` is an intentional performance tradeoff for this model. In architecture reviews, treat this as expected unless there is a requirement to support concurrent multi-run execution in one process.

## Template variable change checklist

When adding or changing template variables, update all of the following in the same change:

- Runtime/template variable source in `src/context/runtime.go` and related config wiring in `src/config/config.go`.
- Template coverage tests in `src/config/init_test.go` and helper/unit tests in `src/config/config_test.go` or `src/shell/template_test.go`.
- Website documentation in `website/docs/setup/templates.mdx`.

## Schema sync requirement

When config structure, fields, defaults, or validation behavior changes, keep schema files in sync in the same change.
This is part of the definition of done.

- Update `website/static/schema.json`.
- Update `src/config/schema.json`.

## PR update requirement

For this fork (`trajano/aliae`), open and update pull requests against `master`.

When updating an open PR, the agent must ensure the PR description still matches the current change scope.
The agent must also add an update comment on the PR summarizing what changed in the update.
If any PR check fails, the agent must reproduce and run the equivalent check locally when possible,
fix the issue, and only push after the local check passes.

## Make a PR command

If I ask the agent to make a PR.  It means the following:

1. Commit all changes make sure tree is clean.  Also make sure the commits are appropriate.
2. Push up a branch (make sure it's not the same name as as origin/HEAD)
3. Create the PR against that branch.
4. Wait (poll every 30 seconds) for the PR checks to pass, if not update the PR accordingly.
5. Complete the PR and delete the source branch off github when merged.
