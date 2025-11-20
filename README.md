# clippy

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/clippy)
[![License](https://img.shields.io/github/license/thediveo/clippy)](https://img.shields.io/github/license/thediveo/clippy)
![build and test](https://github.com/thediveo/clippy/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
![goroutines](https://img.shields.io/badge/go%20routines-not%20leaking-success)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/clippy)](https://goreportcard.com/report/github.com/thediveo/clippy)
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)

...an opinionated modular cobra CLI support Go module, with the emphasis being
on "_modular_". Instead of lumping all flag configuration and processing into
one big mess, why not neatly cutting it into logical chunks (or "plugins").

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/clippy`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `clippy.code-workspace` and off you go...

## Supported Go Versions

`clippy` supports versions of Go that are noted by the [Go release
policy](https://golang.org/doc/devel/release.html#policy), that is, major
versions _N_ and _N_-1 (where _N_ is the current major version).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright and License

`clippy` is Copyright 2024, 2025 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
