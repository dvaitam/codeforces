# Codeforces Solutions

This repository contains my personal experiments and solutions for Codeforces problems. The directory tree mirrors the contest numbers (`0-999`, `1000-1999`, `2000-2999`, ...). Inside each contest folder you will typically find:

- `problemX.txt` – the problem statement or notes for problem X.
- Solution files in Go or other languages (`1234A.go`, `solB.cpp`, etc.).
- Optional verifier programs such as `verifierA.go` that compile a submission and run it against sample tests.

## Utilities

Several helper tools live in the repository root:

- `create.go` – scaffold a new contest/problem directory.
- `auto.go` – convert existing C++ solutions to Go using the Codex CLI.
- `auto_so.go` – generate Go solutions from problem statements with Codex.
- `webserver.go` – simple HTTP server for browsing problems and submitting solutions locally.
- Extra programs like `brute.go` or `probe.go` are experimental helpers.

All Go utilities build with the standard tooling. Example:

```bash
cd 1000-1999/1900-1999/1990-1999/1994
go run verifierA.go ./1994A
```

The repository is mostly experimental and does not contain a complete archive of problems. Feel free to use the code or utilities as a reference for your own workflow.

Some contests include verifier programs with deterministic test cases. For
example, contest 90 provides verifiers for problems A and B:

```bash
cd 0-999/0-99/90-99/90
go run verifierA.go ./mySolutionA
go run verifierB.go ./mySolutionB
```
