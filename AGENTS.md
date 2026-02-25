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

`cd src && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run`

If any Markdown or MDX files changed, run this command from the repo root:

`npx --yes markdownlint-cli $(rg --files -g "*.md" -g "*.mdx")`

If any file under `website/**` changed, run this command from the repo root:

`cd website && npm ci && npm run build`

If any file under `packages/scoop/**` changed, run this command from the repo root:

If multiple scopes changed, run all applicable check sets.

If neither `src/**` nor `website/**` changed, these scoped pre-commit checks are not required.

## Website documentation requirement

When a new feature is added, the agent must also update the website documentation in the same change.
For features added in this fork since `https://github.com/JanDeDobbeleer/aliae`
`main`, update both `README.md` and `website/docs/introduction.mdx` as part of PR definition of done.

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
