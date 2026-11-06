# AGENTS.md

Test: `just test` | Before committing: `just ready`

Read [CONTRIBUTING.md](CONTRIBUTING.md) first. It covers prerequisites, setup,
the project layout, code style, test doubles, documentation, testing and
coverage, and pull requests. All of it applies to agents exactly as it applies
to people. This file carries only what is specific to agents.

## Running tools

Invoke tools through `mise`, not from your path:

```bash
mise exec -- just test
```

`mise` is active in a person's shell and supplies the versions `.mise.toml`
declares. An agent's shell has no activation, so a bare `just` resolves to
whatever is installed globally, usually an older version.

The symptom is a check that fails here and passes in continuous integration, on
a file nobody edited. When that happens, establish which version ran before
treating the failure as real.

Run Psion from the checkout with `mise exec -- just run --help`, not a `psion`
on your path. The recipe copies `resources.d/` into place before compiling, and
an installed binary may be older than the source you are changing.

## Scope

- Do not commit unless the person asks you to.
- A documentation or repository setup change does not authorize moving Go
  packages or redesigning providers.
- Do not apply the example resources to this host to test a documentation or
  styling change. They create, remove and change modes on real paths under
  `/tmp`.
- Run the checks the change calls for, then say which passed and which could not
  run.

## Commit trailer

When committing via Claude Code, end the message with:

```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```
