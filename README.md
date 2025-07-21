# Codeforces Solutions

This repository stores my experiments and solutions for Codeforces problems. The directories are organised by contest and problem ID. Each folder often contains:

- `problemX.txt` – the problem statement or notes.
- Solution implementations (`1994A.go`, `solA.cpp`, etc.).
- Verifier programs written in Go (`verifierA.go`, `verifierB.go`, ...). These compile a user solution and check it against provided test cases.

A few utility programs exist in the repository root:

- `create.go` – helper for creating new contest folders and boilerplate files.
- `auto.go` – simple concurrent build runner for multiple solutions.
- `webserver.go` – demo HTTP server used for uploading solutions locally.
- Python submissions run with `python3`.

To build or run any of the Go utilities, use the standard Go tooling. For example:

```bash
# Build and run a verifier
cd 1000-1999/1900-1999/1990-1999/1994
go run verifierA.go ./1994A
```

The repository is mostly experimental and does not contain complete problem sets. Feel free to browse the directories for reference implementations or to adapt the utilities for your own workflow.

